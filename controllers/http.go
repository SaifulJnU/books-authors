package controllers

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// //////////for author controller///////////////
type AuthorController struct {
	db     *mongo.Client
	logger *zap.Logger // Add a logger field
}

func NewAuthorController(db *mongo.Client, logger *zap.Logger) *AuthorController {
	return &AuthorController{
		db:     db,
		logger: logger, // Initialize the logger field
	}
}

// //////////for book controller///////////////
type BookController struct {
	db     *mongo.Client
	logger *zap.Logger // Add a logger field
}

func NewBookController(db *mongo.Client, logger *zap.Logger) *BookController {
	return &BookController{
		db:     db,
		logger: logger, // Initialize the logger field
	}
}

// ----------------------------------------------------------------

type AuthController struct {
	db     *mongo.Client
	logger *zap.Logger // Add a logger field
}

func NewAuthController(db *mongo.Client, logger *zap.Logger) *AuthController {
	return &AuthController{
		db:     db,
		logger: logger, // Initialize the logger field
	}
}
