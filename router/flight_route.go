package router

import (
	"airplane_booking_go/config"
	"airplane_booking_go/controllers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func FlightRoutes(r *gin.Engine, client *mongo.Client, db string) {
	flightCollection := config.GetCollection(client, db, "flights")
	flightController := controllers.NewFlightController(flightCollection)

	r.POST("/flights", flightController.CreateFlight)
	r.GET("/flights", flightController.GetAllFlights)
	r.GET("/flights/:id", flightController.GetFlightByID)
	r.PUT("/flights/:id", flightController.UpdateFlight)
}
