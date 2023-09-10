// controllers/auth_controller.go

package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saifujnu/books-authors/auth"
	"github.com/saifujnu/books-authors/models"
	"go.mongodb.org/mongo-driver/bson"
)

// Signup handles user registration
func (ac *AuthController) Signup(c *gin.Context) {
	// Parse the user data from the request
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash and salt the user's password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword

	// Store the user in the MongoDB collection
	userCollection := ac.db.Database("book-authors").Collection("users")
	_, err = userCollection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.Status(http.StatusCreated)
}

// Login handles user login and JWT token generation
func (ac *AuthController) Login(c *gin.Context) {
	// Parse the user's login credentials from the request
	var loginData models.User
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the user by username
	userCollection := ac.db.Database("book-authors").Collection("users")
	var user models.User
	err := userCollection.FindOne(context.Background(), bson.M{"username": loginData.Username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify the credentials with the stored hashed password
	if !VerifyPassword(loginData.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate a JWT token using the GenerateToken function
	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Respond with the token
	c.JSON(http.StatusOK, gin.H{"token": token})
}
