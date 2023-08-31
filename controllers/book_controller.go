package controllers

import (
	"context"
	"net/http"

	"github.com/saifujnu/books-authors/db"
	"github.com/saifujnu/books-authors/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookCollection := db.Client.Database("book-authors").Collection("book")
	insertResult, err := bookCollection.InsertOne(context.Background(), book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	// Getting the inserted ID
	bookID := insertResult.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusCreated, gin.H{"_id": bookID})
}

func GetBooks(c *gin.Context) {
	bookCollection := db.Client.Database("book-authors").Collection("book")
	cursor, err := bookCollection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}
	defer cursor.Close(context.Background())

	var books []models.Book
	if err := cursor.All(context.Background(), &books); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode books"})
		return
	}

	c.JSON(http.StatusOK, books)
}

//rest code will be here

func GetBookByID(c *gin.Context) {

}

func UpdateBook(c *gin.Context) {

}

func DeleteBook(c *gin.Context) {

}
