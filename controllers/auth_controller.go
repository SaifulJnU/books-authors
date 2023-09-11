package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saifujnu/books-authors/auth"
	"github.com/saifujnu/books-authors/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap" // Import the Zap logger package
)

// Signup handles user registration
func (ac *AuthController) Signup(c *gin.Context) {
	// Log the start of user registration.
	ac.logger.Info("User registration started")

	// Parse the user data from the request
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		// Log the error and return a bad request response.
		ac.logger.Error("Invalid JSON input", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash and salt the user's password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		// Log the error and return an internal server error response.
		ac.logger.Error("Failed to hash password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword

	// Store the user in the MongoDB collection
	userCollection := ac.db.Database("book-authors").Collection("users")
	_, err = userCollection.InsertOne(context.Background(), user)
	if err != nil {
		// Log the error and return an internal server error response.
		ac.logger.Error("Failed to create user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Log the successful user registration.
	ac.logger.Debug("User registration completed")

	c.Status(http.StatusCreated)
}

// Login handles user login and JWT token generation
func (ac *AuthController) Login(c *gin.Context) {
	// Log the start of user login.
	ac.logger.Info("User login started")

	// Parse the user's login credentials from the request
	var loginData models.User
	if err := c.ShouldBindJSON(&loginData); err != nil {
		// Log the error and return a bad request response.
		ac.logger.Error("Invalid JSON input", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the user by username
	userCollection := ac.db.Database("book-authors").Collection("users")
	var user models.User
	err := userCollection.FindOne(context.Background(), bson.M{"username": loginData.Username}).Decode(&user)
	if err != nil {
		// Log the error and return an unauthorized response.
		ac.logger.Error("Invalid credentials", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify the credentials with the stored hashed password
	if !VerifyPassword(loginData.Password, user.Password) {
		// Log the error and return an unauthorized response.
		ac.logger.Error("Invalid credentials")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate a JWT token using the GenerateToken function
	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		// Log the error and return an internal server error response.
		ac.logger.Error("Failed to generate token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Log the successful user login.
	ac.logger.Debug("User login completed", zap.String("Username", user.Username))

	// Respond with the token
	c.JSON(http.StatusOK, gin.H{"token": token})
}
