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

func (ac *AuthorController) CreateAuthor(c *gin.Context) {
	// Log the start of author creation.
	ac.logger.Info("Creating author")

	// Parse the JSON request body into an Author struct.
	var author models.Author
	if err := c.ShouldBindJSON(&author); err != nil {
		// Log the error and return a bad request response.
		ac.logger.Error("Invalid JSON input", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Access the author collection in MongoDB.
	authorCollection := ac.db.Database("book-authors").Collection("author")

	// Insert the author document into the collection.
	insertResult, err := authorCollection.InsertOne(context.Background(), author)
	if err != nil {
		// Log the error and return an internal server error response.
		ac.logger.Error("Failed to create author", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create author"})
		return
	}

	// Log the successful creation of the author.
	ac.logger.Debug("Author created successfully", zap.String("AuthorID", insertResult.InsertedID.(primitive.ObjectID).Hex()))

	// Include the created author's information in the response.
	c.JSON(http.StatusCreated, gin.H{
		"message": "Author created successfully",
		"author":  author,
	})
}

func (ac *AuthorController) GetAuthors(c *gin.Context) {
	ac.logger.Info("Fetching authors")
	authorCollection := ac.db.Database("book-authors").Collection("author")
	cursor, err := authorCollection.Find(context.Background(), bson.M{})
	if err != nil {
		ac.logger.Error("Failed to fetch authors", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch authors"})
		return
	}
	defer cursor.Close(context.Background())

	var authors []models.Author
	if err := cursor.All(context.Background(), &authors); err != nil {
		ac.logger.Error("Failed to decode authors", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode authors"})
		return
	}

	ac.logger.Debug("Authors fetched successfully")
	c.JSON(http.StatusOK, authors)
}

func (ac *AuthorController) GetAuthorByID(c *gin.Context) {
	authorID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		ac.logger.Error("Invalid author ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	ac.logger.Info("Fetching author by ID", zap.String("AuthorID", authorID))
	authorCollection := ac.db.Database("book-authors").Collection("author")
	var author models.Author
	err = authorCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&author)
	if err != nil {
		ac.logger.Error("Author not found", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		return
	}

	ac.logger.Debug("Author fetched successfully", zap.String("AuthorID", authorID))
	c.JSON(http.StatusOK, author)
}

func (ac *AuthorController) UpdateAuthor(c *gin.Context) {
	authorID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		ac.logger.Error("Invalid author ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	ac.logger.Info("Updating author", zap.String("AuthorID", authorID))
	var updatedAuthor models.Author
	if err := c.ShouldBindJSON(&updatedAuthor); err != nil {
		ac.logger.Error("Invalid JSON input", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authorCollection := ac.db.Database("book-authors").Collection("author")
	update := bson.M{"$set": updatedAuthor}
	_, err = authorCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		ac.logger.Error("Failed to update author", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update author"})
		return
	}

	ac.logger.Debug("Author updated successfully", zap.String("AuthorID", authorID))
	c.Status(http.StatusOK)
}

func (ac *AuthorController) DeleteAuthor(c *gin.Context) {
	authorID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		ac.logger.Error("Invalid author ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	ac.logger.Info("Deleting author", zap.String("AuthorID", authorID))
	authorCollection := ac.db.Database("book-authors").Collection("author")
	_, err = authorCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		ac.logger.Error("Failed to delete author", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete author"})
		return
	}

	ac.logger.Debug("Author deleted successfully", zap.String("AuthorID", authorID))
	c.Status(http.StatusNoContent)
}

// package controllers

// import (
// 	"context"
// 	"net/http"

// 	"github.com/saifujnu/books-authors/models"

// 	"github.com/gin-gonic/gin"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.uber.org/zap"
// )

// func (ac *AuthorController) CreateAuthor(c *gin.Context) {
// 	ac.logger.Info("Creating author")
// 	var author models.Author
// 	if err := c.ShouldBindJSON(&author); err != nil {
// 		ac.logger.Error("Invalid JSON input", zap.Error(err))
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	authorCollection := ac.db.Database("book-authors").Collection("author")
// 	_, err := authorCollection.InsertOne(context.Background(), author)
// 	if err != nil {
// 		ac.logger.Error("Failed to create author", zap.Error(err))
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create author"})
// 		return
// 	}

// 	ac.logger.Debug("Author created successfully")
// 	c.Status(http.StatusCreated)
// }
