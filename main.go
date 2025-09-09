package main

import (
	"airplane_booking_go/config"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
  connectionString := os.Getenv("connectionString")
  config.ConnectDB(connectionString)

}