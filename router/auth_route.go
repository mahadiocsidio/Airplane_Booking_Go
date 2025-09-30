package router

import (
	"airplane_booking_go/config"
	"airplane_booking_go/controllers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRoutes(r *gin.Engine, client *mongo.Client, db string){
	userCollection := config.GetCollection(client, db, "users")
	userController := controllers.NewUserController(userCollection)

	r.POST("/register", userController.Register)
	r.POST("/login", userController.Login)
}