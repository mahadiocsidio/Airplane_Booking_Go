package router

import (
	"airplane_booking_go/config"
	"airplane_booking_go/controllers"
	"airplane_booking_go/middlewares"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func BookRoutes(r *gin.Engine, client *mongo.Client, db string) {
	bookingCollection := config.GetCollection(client, db, "booking")
	flightCollection := config.GetCollection(client, db, "flights")
	bookingController := controllers.NewBookingController(bookingCollection, flightCollection)

	r.POST("/booking",bookingController.CreateBooking )
	r.GET("/booking", bookingController.GetAllBookings)
	r.GET("/booking/user", bookingController.GetUserBookings)
	r.PUT("/booking/:id/cancel", bookingController.CancelBooking, middlewares.AuthMiddleware())
}
