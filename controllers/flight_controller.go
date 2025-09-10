package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"airplane_booking_go/models"
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

	newFlight := models.Flight{
		ID:            primitive.NewObjectID(),
		Airline:       req.Airline,
		FlightNumber:  req.FlightNumber,
		Departure:     req.Departure,
		Arrival:       req.Arrival,
		DepartureTime: req.DepartureTime,
		ArrivalTime:   req.ArrivalTime,
		Duration:      req.Duration,
		Price:         req.Price,
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