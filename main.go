package main

import (
	"github.com/joho/godotenv"
	"github.com/saifujnu/books-authors/controllers"
	"github.com/saifujnu/books-authors/db/mongo"

	"github.com/gin-gonic/gin"
	"github.com/saifujnu/books-authors/config"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	config.SetEnvionment()

}

func main() {
	//task1: take from env
	// err := db.ConnectMongoDB("mongodb://admin:secret@localhost:27017,localhost:27018,localhost:27019/?replicaSet=rs0")
	// if err != nil {
	// 	panic(err)
	// }

	m, err := mongo.Connect()

	if err != nil {
		panic(err)
	}

	router := gin.Default()

	// mgoConn, err := db.GetMongoCon("mongodb://admin:secret@localhost:27017,localhost:27018,localhost:27019/?replicaSet=rs0")
	// if err != nil {
	// 	panic(err)
	// }

	//authorController := controllers.NewAuthorController(mgoConn)
	authorController := controllers.NewAuthorController(m)
	bookController := controllers.NewBookController(m)

	// Author routes
	router.POST("/authors", authorController.CreateAuthor)
	router.GET("/authors", authorController.GetAuthors)
	router.GET("/authors/:id", authorController.GetAuthorByID)
	router.PUT("/authors/:id", authorController.UpdateAuthor)
	router.DELETE("/authors/:id", authorController.DeleteAuthor)

	// Book routes
	router.POST("/books", bookController.CreateBook)
	router.GET("/books", bookController.GetBooks)
	router.GET("/books/:id", bookController.GetBookByID)
	router.PUT("/books/:id", bookController.UpdateBook)
	router.DELETE("/books/:id", bookController.DeleteBook)

	router.Run(":8080")
}
