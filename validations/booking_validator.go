package validations

import (
	// "airplane_booking_go/models"
	// "time"
)

type CreateBookingRequest struct {
	FlightID    string   `json:"flightId" binding:"required"`
	SeatNumbers []string `json:"seatNumbers" binding:"required"`
}

type GetUserBookingsRequest struct {
	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=10"`
}