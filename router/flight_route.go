package router

import (
	"airplane_booking_go/config"
	"airplane_booking_go/controllers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func FlightRoutes(r *gin.Engine, client *mongo.Client, db string){
	flightCollection := config.GetCollection(client, db, "flights")
	flightController := controllers.NewFlightController(flightCollection)

	r.POST("/register", flightController.CreateFlight)
	r.GET("/login", flightController.GetFlights)
}