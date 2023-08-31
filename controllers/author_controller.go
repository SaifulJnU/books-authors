package controllers

import (
	"context"
	"net/http"

	"github.com/saifujnu/books-authors/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// task 2: when you create author then return it
func (a *AuthorController) CreateAuthor(c *gin.Context) {
	var author models.Author
	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authorCollection := a.db.Database("book-authors").Collection("author")
	_, err := authorCollection.InsertOne(context.Background(), author)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create author"})
		return
	}

	c.Status(http.StatusCreated)
}

func (a *AuthorController) GetAuthors(c *gin.Context) {
	authorCollection := a.db.Database("book-authors").Collection("author")
	cursor, err := authorCollection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch authors"})
		return
	}
	defer cursor.Close(context.Background())

	var authors []models.Author
	if err := cursor.All(context.Background(), &authors); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode authors"})
		return
	}

	c.JSON(http.StatusOK, authors)
}

func (a *AuthorController) GetAuthorByID(c *gin.Context) {
	authorID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	authorCollection := a.db.Database("book-authors").Collection("author")
	var author models.Author
	err = authorCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&author)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		return
	}

	c.JSON(http.StatusOK, author)
}

func (a *AuthorController) UpdateAuthor(c *gin.Context) {
	authorID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	var updatedAuthor models.Author
	if err := c.ShouldBindJSON(&updatedAuthor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authorCollection := a.db.Database("book-authors").Collection("author")
	update := bson.M{"$set": updatedAuthor}
	_, err = authorCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update author"})
		return
	}

	c.Status(http.StatusOK)
}

func (a *AuthorController) DeleteAuthor(c *gin.Context) {
	authorID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	authorCollection := a.db.Database("book-authors").Collection("author")
	_, err = authorCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete author"})
		return
	}

	c.Status(http.StatusNoContent)
}
