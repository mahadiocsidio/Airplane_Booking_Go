// models/flight.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Flight struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Airline       string             `bson:"airline" json:"airline"`
	FlightNumber  string             `bson:"flightNumber" json:"flightNumber"`
	Departure     Airport            `bson:"departure" json:"departure"`
	Arrival       Airport            `bson:"arrival" json:"arrival"`
	DepartureTime time.Time          `bson:"departureTime" json:"departureTime"`
	ArrivalTime   time.Time          `bson:"arrivalTime" json:"arrivalTime"`
	Duration      int                `bson:"duration" json:"duration"`
	MinPrice      float64            `bson:"minPrice" json:"minPrice"`
	Seats         []Seat             `bson:"seats" json:"seats"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type Seat struct {
	Number      string  `bson:"number" json:"number"`
	Class       string  `bson:"class" json:"class"` 
	IsAvailable bool    `bson:"isAvailable" json:"isAvailable"`
	Price       float64 `bson:"price" json:"price"`
}

type Airport struct {
	Code    string `bson:"code" json:"code"`
	Name    string `bson:"name" json:"name"`
	City    string `bson:"city" json:"city"`
	Country string `bson:"country" json:"country"`
}
