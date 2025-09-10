package main

import (
	"airplane_booking_go/config"
  "airplane_booking_go/router"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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
	router.UserRoutes(r, client, db)
  router.FlightRoutes(r, client, db)
  r.Run(":8080")
}
