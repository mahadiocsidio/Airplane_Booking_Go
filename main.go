package main

import (
	"airplane_booking_go/config"
  	"airplane_booking_go/router"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/cmd/swag"
	"github.com/swaggo/files"
)
// @title Airplane_Booking API
// @version 1.0
// @description This is a backend for airplane booking system.
// @host localhost:8080
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connectionString := os.Getenv("connectionString")
	db := os.Getenv("db")
	client := config.ConnectDB(connectionString)

	//router setup
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.UserRoutes(r, client, db)
  	router.FlightRoutes(r, client, db)
	router.BookRoutes(r, client, db)
  	r.Run(":8080")
}
