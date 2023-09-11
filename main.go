package main

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	ginzap "github.com/gin-contrib/zap"
	"github.com/saifujnu/books-authors/auth"
	"github.com/saifujnu/books-authors/config"
	"github.com/saifujnu/books-authors/controllers"
	"github.com/saifujnu/books-authors/db/mongo"
)

var Logger *zap.Logger

func InitializeLogger() (*zap.Logger, error) {
	logger, err := zap.NewDevelopment() // with NewDevelopment() I can Call info, error and debug also
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	config.SetEnvionment()

	// Replace the standard logger with Zap logger
	var errLogger error
	Logger, errLogger = InitializeLogger()
	if errLogger != nil {
		panic("Failed to initialize Zap logger: " + errLogger.Error())
	}
}

func main() {
	// Initialize MongoDB connection
	m, err := mongo.Connect()
	if err != nil {
		Logger.Error("Failed to connect to MongoDB", zap.Error(err))
		os.Exit(1)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	//combining Gin and Zap logger
	router.Use(ginzap.Ginzap(Logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(Logger, true))

	// Initialize controllers with the MongoDB client
	authorController := controllers.NewAuthorController(m, Logger)
	bookController := controllers.NewBookController(m, Logger)
	authController := controllers.NewAuthController(m, Logger)

	// Setup authentication routes
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/signup", authController.Signup)
		authRoutes.POST("/login", authController.Login)
	}

	// Define your CRUD routes for books and authors here
	bookRoutes := router.Group("/books")
	bookRoutes.Use(auth.JWTMiddleware()) // Protect /books routes with JWT middleware
	{
		bookRoutes.GET("/", bookController.GetBooks)
		bookRoutes.GET("/:id", bookController.GetBookByID)
		bookRoutes.POST("/", bookController.CreateBook)
		bookRoutes.PUT("/:id", bookController.UpdateBook)
		bookRoutes.DELETE("/:id", bookController.DeleteBook)
		bookRoutes.GET("/books-and-authors", bookController.GetAllBooksAndAuthors)
		bookRoutes.GET("/books-by-author/:authorName", bookController.GetBooksByAuthorName)
	}

	authorRoutes := router.Group("/authors")
	authorRoutes.Use(auth.JWTMiddleware()) // Protect /authors routes with JWT middleware
	{
		authorRoutes.GET("/", authorController.GetAuthors)
		authorRoutes.GET("/:id", authorController.GetAuthorByID)
		authorRoutes.POST("/", authorController.CreateAuthor)
		authorRoutes.PUT("/:id", authorController.UpdateAuthor)
		authorRoutes.DELETE("/:id", authorController.DeleteAuthor)
	}

	Logger.Info("Server started on :8080") // Use the logger

	err1 := router.Run(":8080")
	if err1 != nil {
		Logger.Error("Failed to start server", zap.Error(err1))
		os.Exit(1)
	}

}
