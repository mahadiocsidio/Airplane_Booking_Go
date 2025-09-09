package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	UserID     primitive.ObjectID   `bson:"userId" json:"userId"`
	FlightID   primitive.ObjectID   `bson:"flightId" json:"flightId"`
	Seats      []string             `bson:"seats" json:"seats"` // seat numbers, ex: ["12A", "12B"]
	TotalPrice float64              `bson:"totalPrice" json:"totalPrice"`
	Status     string               `bson:"status" json:"status"` // pending, confirmed, cancelled
	BookedAt   time.Time            `bson:"bookedAt" json:"bookedAt"`
	UpdatedAt  time.Time            `bson:"updatedAt" json:"updatedAt"`
}
