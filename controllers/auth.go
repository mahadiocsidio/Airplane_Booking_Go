package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"airplane_booking_go/models"
	"airplane_booking_go/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	UserCollection *mongo.Collection
}

func NewUserController(userCollection *mongo.Collection) *UserController {
	return &UserController{UserCollection: userCollection}
}

// ========== REQUEST STRUCT ==========
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ========== HANDLERS ==========
// Register a User
func (uc *UserController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	newUser := models.User{
		ID:        primitive.NewObjectID(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword), // hashed password
		Role:      "user",                 // default role
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = uc.UserCollection.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	// jangan return password ke client
	c.JSON(http.StatusCreated, gin.H{
		"message"	: "user registered successfully",
		"code"		: "201",
		"status"	: "Created",
	})
}

// Log in With User Data
func (uc *UserController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// cari user by email
	var user models.User
	err := uc.UserCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// cek password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// TODO: generate JWT token nanti
	token, err := utils.GenerateToken(user.ID.Hex(), user.Role)
	if err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
	return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message"	: "Login Succes",
		"code"		: "200",
		"status"	: "OK",
		"token"		: token})
}
