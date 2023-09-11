package controllers

import (
	"context"
	"net/http"

	"github.com/saifujnu/books-authors/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func (bc *BookController) CreateBook(c *gin.Context) {
	// Log the start of book creation.
	bc.logger.Info("Creating book")

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		// Log the error and return a bad request response.
		bc.logger.Error("Invalid JSON input", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookCollection := bc.db.Database("book-authors").Collection("book")
	insertResult, err := bookCollection.InsertOne(context.Background(), book)
	if err != nil {
		// Log the error and return an internal server error response.
		bc.logger.Error("Failed to create book", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	// Log the successful creation of the book.
	bc.logger.Debug("Book created successfully", zap.String("BookID", insertResult.InsertedID.(primitive.ObjectID).Hex()))

	// Include the created book's information in the response.
	c.JSON(http.StatusCreated, gin.H{
		"message": "Book created successfully",
		"book":    book,
	})
}

func (bc *BookController) GetBooks(c *gin.Context) {
	// Log the start of fetching books.
	bc.logger.Info("Fetching books")

	bookCollection := bc.db.Database("book-authors").Collection("book")
	cursor, err := bookCollection.Find(context.Background(), bson.M{})
	if err != nil {
		// Log the error and return an internal server error response.
		bc.logger.Error("Failed to fetch books", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}
	defer cursor.Close(context.Background())

	var books []models.Book
	if err := cursor.All(context.Background(), &books); err != nil {
		// Log the error and return an internal server error response.
		bc.logger.Error("Failed to decode books", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode books"})
		return
	}

	// Log the successful fetching of books.
	bc.logger.Debug("Books fetched successfully")
	c.JSON(http.StatusOK, books)
}

func (bc *BookController) GetBookByID(c *gin.Context) {
	bookID := c.Param("id")
	bookObjID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		// Log the error and return a bad request response.
		bc.logger.Error("Invalid book ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// Log the start of fetching a book by ID.
	bc.logger.Info("Fetching book by ID", zap.String("BookID", bookID))

	bookCollection := bc.db.Database("book-authors").Collection("book")
	var book models.Book
	err = bookCollection.FindOne(context.Background(), bson.M{"_id": bookObjID}).Decode(&book)
	if err != nil {
		// Log the error and return a not found response.
		bc.logger.Error("Book not found", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// Log the successful fetching of the book.
	bc.logger.Debug("Book fetched successfully", zap.String("BookID", bookID))
	c.JSON(http.StatusOK, book)
}

func (bc *BookController) UpdateBook(c *gin.Context) {
	bookID := c.Param("id")
	bookObjID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		// Log the error and return a bad request response.
		bc.logger.Error("Invalid book ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// Log the start of updating a book.
	bc.logger.Info("Updating book", zap.String("BookID", bookID))

	var updatedBook models.Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		// Log the error and return a bad request response.
		bc.logger.Error("Invalid JSON input", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookCollection := bc.db.Database("book-authors").Collection("book")
	update := bson.M{"$set": updatedBook}
	_, err = bookCollection.UpdateOne(context.Background(), bson.M{"_id": bookObjID}, update)
	if err != nil {
		// Log the error and return an internal server error response.
		bc.logger.Error("Failed to update book", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	// Log the successful update of the book.
	bc.logger.Debug("Book updated successfully", zap.String("BookID", bookID))
	c.Status(http.StatusOK)
}

func (bc *BookController) DeleteBook(c *gin.Context) {
	bookID := c.Param("id")
	bookObjID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		// Log the error and return a bad request response.
		bc.logger.Error("Invalid book ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// Log the start of deleting a book.
	bc.logger.Info("Deleting book", zap.String("BookID", bookID))

	bookCollection := bc.db.Database("book-authors").Collection("book")
	_, err = bookCollection.DeleteOne(context.Background(), bson.M{"_id": bookObjID})
	if err != nil {
		// Log the error and return an internal server error response.
		bc.logger.Error("Failed to delete book", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	// Log the successful deletion of the book.
	bc.logger.Debug("Book deleted successfully", zap.String("BookID", bookID))
	c.Status(http.StatusNoContent)
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
		// Log the error and return an internal server error response.
		bc.logger.Error("Failed to aggregate books and authors", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to aggregate books and authors"})
		return
	}
	defer cursor.Close(context.Background())

	var combinedList []bson.M
	if err := cursor.All(context.Background(), &combinedList); err != nil {
		// Log the error and return an internal server error response.
		bc.logger.Error("Failed to decode books and authors", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode books and authors"})
		return
	}

	// Log the successful aggregation and response.
	bc.logger.Debug("Books and authors aggregated successfully")
	c.JSON(http.StatusOK, combinedList)
}

func (bc *BookController) GetBooksByAuthorName(c *gin.Context) {
	authorName := c.Param("authorName")

	// Find the author by first name
	authorCollection := bc.db.Database("book-authors").Collection("author")
	var author models.Author

	err := authorCollection.FindOne(context.Background(), bson.M{"firstName": authorName}).Decode(&author)
	if err != nil {
		// Log the error and return a not found response.
		bc.logger.Error("Author not found", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		return
	}

	// Find the books by author ID
	bookCollection := bc.db.Database("book-authors").Collection("book")
	cursor, err := bookCollection.Find(context.Background(), bson.M{"authorId": author.ID})
	if err != nil {
		// Log the error and return an internal server error response.
		bc.logger.Error("Failed to fetch books", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}
	defer cursor.Close(context.Background())

	var books []models.Book
	if err := cursor.All(context.Background(), &books); err != nil {
		// Log the error and return an internal server error response.
		bc.logger.Error("Failed to decode books", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode books"})
		return
	}

	// Log the successful response.
	bc.logger.Debug("Books fetched by author name successfully", zap.String("AuthorName", authorName))
	c.JSON(http.StatusOK, books)
}
