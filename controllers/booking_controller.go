package controllers

import(
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"airplane_booking_go/models"
)

type BookingController struct {
	BookingCollection *mongo.Collection
	FlightCollection  *mongo.Collection
}

func NewBookingController(bookingColl, flightColl *mongo.Collection) *BookingController {
	return &BookingController{
		BookingCollection: bookingColl,
		FlightCollection:  flightColl,
	}
}

type CreateBookingRequest struct {
	FlightID    string   `json:"flightId" binding:"required"`
	SeatNumbers []string `json:"seatNumbers" binding:"required"`
}

// GetFlightByID → fetch flight detail by ID
func (fc *FlightController) GetFlightByID(c *gin.Context) {
	flightId := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(flightId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid flight id"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var flight models.Flight
	err = fc.FlightCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&flight)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "flight not found"})
		return
	}

	c.JSON(http.StatusOK, flight)
}

// CreateBooking → user booking kursi
func (bc *BookingController) CreateBooking(c *gin.Context) {
	var req CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// fetch userID from context
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	flightObjID, err := primitive.ObjectIDFromHex(req.FlightID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid flightId"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// fetch flight data
	var flight models.Flight
	if err := bc.FlightCollection.FindOne(ctx, bson.M{"_id": flightObjID}).Decode(&flight); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "flight not found"})
		return
	}

	// check seat availability & calculate total
	var selectedSeats []models.Seat
	totalPrice := 0.0
	for _, seatNum := range req.SeatNumbers {
		found := false
		for _, seat := range flight.Seats {
			if seat.Number == seatNum {
				found = true
				if !seat.IsAvailable {
					c.JSON(http.StatusBadRequest, gin.H{"error": "seat " + seatNum + " is not available"})
					return
				}
				selectedSeats = append(selectedSeats, seat)
				totalPrice += seat.Price
				break
			}
		}
		if !found {
			c.JSON(http.StatusBadRequest, gin.H{"error": "seat " + seatNum + " not found"})
			return
		}
	}

	// update seat availability
	for _, seatNum := range req.SeatNumbers {
		result, err := bc.FlightCollection.UpdateOne(
			ctx,
			bson.M{
				"_id": flightObjID,
				"seats": bson.M{
					"$elemMatch": bson.M{"number": seatNum, "is_available": true},
				},
			},
			bson.M{
				"$set": bson.M{
					"seats.$.is_available": false,
					"updated_at":           time.Now(),
				},
			},
		)
		if err != nil || result.ModifiedCount == 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "seat " + seatNum + " just got booked"})
			return
		}
	}

	// create booking
	booking := models.Booking{
		ID:         primitive.NewObjectID(),
		UserID:     userID.(primitive.ObjectID),
		FlightID:   flightObjID,
		Seats:      selectedSeats,
		TotalPrice: totalPrice,
		Status:     "confirmed",
		BookedAt:   time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err = bc.BookingCollection.InsertOne(ctx, booking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create booking"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "booking created", "booking": booking})
}

// GetBookings → fetch all booking user
func (bc *BookingController) GetBookings(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := bc.BookingCollection.Find(ctx, bson.M{"userId": userID.(primitive.ObjectID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookings"})
		return
	}
	defer cursor.Close(ctx)

	var bookings []models.Booking
	if err = cursor.All(ctx, &bookings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// update booking to cancelled
func (bc *BookingController) CancelBooking(c *gin.Context) {
	bookingId := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(bookingId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var booking models.Booking
	if err := bc.BookingCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&booking); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	// update booking status
	_, err = bc.BookingCollection.UpdateOne(ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"status": "canceled", "updated_at": time.Now()}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel booking"})
		return
	}

	// update seats to avaiable
	bc.FlightCollection.UpdateOne(ctx,
		bson.M{"_id": booking.FlightID, "seats": booking.Seats},
		bson.M{"$set": bson.M{"seats.$.is_available": true}},
	)

	c.JSON(http.StatusOK, gin.H{"message": "booking canceled successfully"})
}
