package main

import (
	"airplane_booking_go/config"
	"airplane_booking_go/controllers"
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
	userCollection := config.GetCollection(client, db, "users")

	//controller init
	userController := controllers.NewUserController(userCollection)

	//router setup
	r := gin.Default()
	r.POST("/register", userController.Register)
	r.GET("/login", userController.Login)
	r.Run(":8080")
}
