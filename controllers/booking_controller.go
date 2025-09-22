package controllers

import(
	"context"
	"net/http"
	"fmt"
	"time"
	"errors"

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

	ctx := context.Background()
	session, err := bc.FlightCollection.Database().Client().StartSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start session"})
		return
	}
	defer session.EndSession(ctx)

	var booking models.Booking

	// transaction function
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// fetch flight data
		var flight models.Flight
		if err := bc.FlightCollection.FindOne(sessCtx, bson.M{"_id": flightObjID}).Decode(&flight); err != nil {
			return nil, fmt.Errorf("flight not found")
		}

		// check seats avaiable + count total
		var selectedSeats []models.Seat
		totalPrice := 0.0
		for _, seatNum := range req.SeatNumbers {
			found := false
			for _, seat := range flight.Seats {
				if seat.Number == seatNum {
					found = true
					if !seat.IsAvailable {
						return nil, fmt.Errorf("seat %s not available", seatNum)
					}
					selectedSeats = append(selectedSeats, seat)
					totalPrice += seat.Price
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("seat %s not found", seatNum)
			}
		}

		// update seats into unavailable (bulk update)
		for _, seatNum := range req.SeatNumbers {
			result, err := bc.FlightCollection.UpdateOne(
				sessCtx,
				bson.M{
					"_id": flightObjID,
					"seats": bson.M{
						"$elemMatch": bson.M{"number": seatNum, "isAvailable": true},
					},
				},
				bson.M{
					"$set": bson.M{
						"seats.$.isAvailable": false,
						"updated_at":           time.Now(),
					},
				},
			)
			if err != nil || result.ModifiedCount == 0 {
				return nil, fmt.Errorf("seat %s just got booked", seatNum)
			}
		}

		// buat booking baru
		booking = models.Booking{
			ID:         primitive.NewObjectID(),
			UserID:     userID.(primitive.ObjectID),
			FlightID:   flightObjID,
			Seats:      selectedSeats,
			TotalPrice: totalPrice,
			Status:     "confirmed", // nanti bisa diganti "pending" kalau ada pembayaran
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if _, err := bc.BookingCollection.InsertOne(sessCtx, booking); err != nil {
			return nil, fmt.Errorf("failed to insert booking: %v", err)
		}

		return nil, nil
	}

	// jalankan transaksi
	if _, err = session.WithTransaction(ctx, callback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "booking created",
		"booking": booking,
	})
}

// GetBookings → fetch all booking user
func (bc *BookingController) GetAllBookings(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// filter opsional
	filter := bson.M{}
	if status := c.Query("status"); status != "" {
		filter["status"] = status
	}

	cursor, err := bc.BookingCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookings"})
		return
	}
	defer cursor.Close(ctx)

	var bookings []models.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse bookings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "OK",
		"bookings": bookings,
	})
}

func (bc *BookingController) GetUserBookings(c *gin.Context) {
	// ambil userId dari context (udah di-set waktu auth)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := bc.BookingCollection.Find(ctx, bson.M{
		"userId": userID.(primitive.ObjectID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookings"})
		return
	}
	defer cursor.Close(ctx)

	var bookings []models.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse bookings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "OK",
		"bookings": bookings,
	})
}

// update booking to cancelled
func (bc *BookingController) CancelBooking(c *gin.Context) {
	bookingID := c.Param("id")
	bookingObjID, err := primitive.ObjectIDFromHex(bookingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bookingId"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// start session
	session, err := bc.BookingCollection.Database().Client().StartSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start session"})
		return
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		// fetch booking
		var booking models.Booking
		if err := bc.BookingCollection.FindOne(sc, bson.M{"_id": bookingObjID}).Decode(&booking); err != nil {
			return err
		}

		if booking.Status != "confirmed" {
			return errors.New("booking is not active")
		}

		// update booking status
		_, err = bc.BookingCollection.UpdateOne(
			sc,
			bson.M{"_id": bookingObjID},
			bson.M{"$set": bson.M{"status": "cancelled", "updated_at": time.Now()}},
		)
		if err != nil {
			return err
		}

		// release seats back to available
		for _, seat := range booking.Seats {
			_, err := bc.FlightCollection.UpdateOne(
				sc,
				bson.M{"_id": booking.FlightID, "seats.number": seat.Number},
				bson.M{"$set": bson.M{"seats.$.is_available": true, "updated_at": time.Now()}},
			)
			if err != nil {
				return err
			}
		}

		return session.CommitTransaction(sc)
	})

	if err != nil {
		session.AbortTransaction(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "booking cancelled successfully"})
}

