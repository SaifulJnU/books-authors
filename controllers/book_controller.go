package controllers

import (
	"context"
	"net/http"

	"github.com/saifujnu/books-authors/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (bc *BookController) CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookCollection := bc.db.Database("book-authors").Collection("book")
	insertResult, err := bookCollection.InsertOne(context.Background(), book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	// Getting the inserted ID
	bookID := insertResult.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusCreated, gin.H{"_id": bookID})

}

func (bc *BookController) GetBooks(c *gin.Context) {
	bookCollection := bc.db.Database("book-authors").Collection("book")
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

func (bc *BookController) GetBookByID(c *gin.Context) {
	bookID := c.Param("id")
	bookObjID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	bookCollection := bc.db.Database("book-authors").Collection("book")
	var book models.Book

	err = bookCollection.FindOne(context.Background(), bson.M{"_id": bookObjID}).Decode(&book)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (bc *BookController) UpdateBook(c *gin.Context) {
	bookID := c.Param("id")
	bookObjID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var updatedBook models.Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookCollection := bc.db.Database("book-authors").Collection("book")

	updateResult, err := bookCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": bookObjID},
		bson.D{{Key: "$set", Value: updatedBook}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	if updateResult.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book updated successfully"})
}

func (bc *BookController) DeleteBook(c *gin.Context) {
	bookID := c.Param("id")
	bookObjID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	bookCollection := bc.db.Database("book-authors").Collection("book")

	deleteResult, err := bookCollection.DeleteOne(context.Background(), bson.M{"_id": bookObjID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	if deleteResult.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

// combine controller
func (bc *BookController) GetAllBooksAndAuthors(c *gin.Context) {
	bookCollection := bc.db.Database("book-authors").Collection("book")
	//authorCollection := bc.db.Database("book-authors").Collection("author")

	// Define the aggregation pipeline
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "author",
				"localField":   "authorId",
				"foreignField": "_id",
				"as":           "authorInfo",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$authorInfo",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$project": bson.M{
				"_id":        1,
				"title":      1,
				"authorInfo": 1,
			},
		},
	}

	cursor, err := bookCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to aggregate books and authors"})
		return
	}
	defer cursor.Close(context.Background())

	var combinedList []bson.M
	if err := cursor.All(context.Background(), &combinedList); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode books and authors"})
		return
	}

	c.JSON(http.StatusOK, combinedList)
}

func (bc *BookController) GetBooksByAuthorName(c *gin.Context) {
	authorName := c.Param("authorName")

	// Find the author by full name
	authorCollection := bc.db.Database("book-authors").Collection("author")
	var author models.Author

	err := authorCollection.FindOne(context.Background(), bson.M{"fullName": authorName}).Decode(&author)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		return
	}

	// Find the books by author ID
	bookCollection := bc.db.Database("book-authors").Collection("book")
	cursor, err := bookCollection.Find(context.Background(), bson.M{"authorId": author.ID})
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
