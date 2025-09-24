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

	booking := r.Group("/booking", middlewares.AuthMiddleware())
	{
    	booking.POST("/", middlewares.AuthMiddleware(), bookingController.CreateBooking)
    	booking.GET("/", bookingController.GetAllBookings)
    	booking.GET("/user", bookingController.GetUserBookings)
    	booking.PUT("/:id/cancel", bookingController.CancelBooking)
	}
}
