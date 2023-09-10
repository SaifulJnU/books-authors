package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/saifujnu/books-authors/auth"
	"github.com/saifujnu/books-authors/controllers" // Import the controllers package
)

// SetupAuthRoutes sets up authentication routes
func SetupAuthRoutes(router *gin.Engine) {
	authController := controllers.NewAuthController(nil) // Initialize your AuthController

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/signup", authController.Signup)
		authRoutes.POST("/login", authController.Login)
	}

	// Apply JWT middleware to protected routes
	authRoutes.Use(auth.JWTMiddleware())

	// Define your CRUD routes for books and authors here
}
