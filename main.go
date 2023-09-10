// main.go

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/saifujnu/books-authors/auth"
	"github.com/saifujnu/books-authors/config"
	"github.com/saifujnu/books-authors/controllers"
	"github.com/saifujnu/books-authors/db/mongo"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	config.SetEnvionment()

}

func main() {
	// Initialize MongoDB connection
	// mongoURI := "mongodb://admin:secret@localhost:27017,localhost:27018,localhost:27019/?replicaSet=rs0" // Replace with your MongoDB URI
	// clientOptions := options.Client().ApplyURI(mongoURI)
	// client, err := mongo.Connect(context.Background(), clientOptions)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer client.Disconnect(context.Background())

	m, err := mongo.Connect()

	if err != nil {
		panic(err)
	}
	router := gin.Default()

	// Initialize controllers with the MongoDB client
	authorController := controllers.NewAuthorController(m)
	bookController := controllers.NewBookController(m)
	authController := controllers.NewAuthController(m)

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

	// Start the server
	// ...
	router.Run(":8080")
}
