package validations

import(
	"time"
	// "github.com/go-playground/validator/v10"
	"airplane_booking_go/models"
)

// ========== REQUEST STRUCT ==========
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

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