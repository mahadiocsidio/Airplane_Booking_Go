package router

import (
	"airplane_booking_go/config"
	"airplane_booking_go/controllers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func BookRoutes(r *gin.Engine, client *mongo.Client, db string) {
	bookingCollection := config.GetCollection(client, db, "booking")
	flightCollection := config.GetCollection(client, db, "flights")
	bookingController := controllers.NewBookingController(bookingCollection, flightCollection)

	r.POST("/flights",bookingController.CreateBooking )
	r.GET("/flights", bookingController.GetBookings)
	r.PUT("/flights/:id/cancel", bookingController.GetBookings)
}
