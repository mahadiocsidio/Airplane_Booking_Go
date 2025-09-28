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
    	booking.POST("/book", middlewares.AuthMiddleware(), bookingController.CreateBooking)
    	booking.GET("/book", bookingController.GetAllBookings)
    	booking.GET("/user/book", middlewares.AuthMiddleware(), bookingController.GetUserBookings)
    	booking.GET("/book/:id", middlewares.AuthMiddleware(), bookingController.GetUserBookingDetail)
    	booking.PUT("/book/:id/cancel", middlewares.AuthMiddleware(), bookingController.CancelBooking)
	}
}
