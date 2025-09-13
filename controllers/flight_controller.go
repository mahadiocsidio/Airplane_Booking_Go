package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"


	"airplane_booking_go/models"
	"airplane_booking_go/validations"
	"airplane_booking_go/utils"
)

type FlightController struct {
	FlightCollection *mongo.Collection
}

func NewFlightController(flightCollection *mongo.Collection) *FlightController {
	return &FlightController{FlightCollection: flightCollection}
}

// ========== REQUEST STRUCT ==========
type CreateFlightRequest struct {
	Airline       string          `json:"airline" binding:"required"`
	FlightNumber  string          `json:"flightNumber" binding:"required"`
	Departure     models.Airport  `json:"departure" binding:"required"`
	Arrival       models.Airport  `json:"arrival" binding:"required"`
	DepartureTime time.Time       `json:"departureTime" binding:"required"`
	ArrivalTime   time.Time       `json:"arrivalTime" binding:"required"`
	Duration      int             `json:"duration" binding:"required"`
	Price         float64         `json:"price" binding:"required"`
	Seats         []models.Seat   `json:"seats" binding:"required"`
}

// ========== HANDLERS ==========

// Create flight (admin only ideally)
func (fc *FlightController) CreateFlight(c *gin.Context) {
	var req CreateFlightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// searching for cheapest seat price
	minPrice := req.Seats[0].Price
	for _, seat := range req.Seats {
		if seat.Price < minPrice {
			minPrice = seat.Price
		}
	}

	newFlight := models.Flight{
		ID:            primitive.NewObjectID(),
		Airline:       req.Airline,
		FlightNumber:  req.FlightNumber,
		Departure:     req.Departure,
		Arrival:       req.Arrival,
		DepartureTime: req.DepartureTime,
		ArrivalTime:   req.ArrivalTime,
		Duration:      req.Duration,
		Price:         minPrice, // start from
		Seats:         req.Seats,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := fc.FlightCollection.InsertOne(ctx, newFlight)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert data"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code"		:"200",
		"status"	:"OK",
		"message"	:"data created", 
		"flight"	: newFlight})
}

// Get all flights
func (fc *FlightController) GetFlights(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := fc.FlightCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch flights"})
		return
	}
	defer cursor.Close(ctx)

	var flights []models.Flight
	if err = cursor.All(ctx, &flights); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode flights"})
		return
	}

	c.JSON(http.StatusOK, flights)
}

// UpdateFlight → update flight data (admin only ideally)
func (fc *FlightController) UpdateFlight(c *gin.Context) {
	flightId := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(flightId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid flight id"})
		return
	}

	var req CreateFlightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hitung harga seat termurah
	if len(req.Seats) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "seats cannot be empty"})
		return
	}
	minPrice := req.Seats[0].Price
	for _, seat := range req.Seats {
		if seat.Price < minPrice {
			minPrice = seat.Price
		}
	}

	update := bson.M{
		"airline":        req.Airline,
		"flight_number":  req.FlightNumber,
		"departure":      req.Departure,
		"arrival":        req.Arrival,
		"departure_time": req.DepartureTime,
		"arrival_time":   req.ArrivalTime,
		"duration":       req.Duration,
		"price":          minPrice,
		"seats":          req.Seats,
		"updated_at":     time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := fc.FlightCollection.UpdateOne(ctx,
		bson.M{"_id": objID},
		bson.M{"$set": update},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update flight"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "flight not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "flight updated successfully"})
}

func (fc *FlightController) SearchFlights(c *gin.Context) {
	var req validations.SearchFlightRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p := utils.GetPagination(c)

	// build filter
	filter := bson.M{}
	if req.From != "" {
		filter["departure.code"] = req.From
	}
	if req.To != "" {
		filter["arrival.code"] = req.To
	}
	if req.Date != "" {
		day, err := time.Parse("2006-01-02", req.Date)
		if err == nil {
			start := day
			end := day.Add(24 * time.Hour)
			filter["departure_time"] = bson.M{
				"$gte": start,
				"$lt":  end,
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// hitung totalCount untuk pagination
	totalCount, err := fc.FlightCollection.CountDocuments(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count flights"})
		return
	}

	opts := options.Find().
		SetSkip(int64(p.Skip)).
		SetLimit(int64(p.Limit)).
		SetSort(bson.M{"departure_time": 1})

	cursor, err := fc.FlightCollection.Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch flights"})
		return
	}
	defer cursor.Close(ctx)

	var flights []models.Flight
	if err = cursor.All(ctx, &flights); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode flights"})
		return
	}

	totalPages := int((totalCount + int64(p.Limit) - 1) / int64(p.Limit)) // ceil

	c.JSON(http.StatusOK, gin.H{
		"page":        p.Page,
		"limit":       p.Limit,
		"count":       len(flights),
		"totalCount":  totalCount,
		"totalPages":  totalPages,
		"flights":     flights,
	})
}