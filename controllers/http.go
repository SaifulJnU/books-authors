package controllers

import "go.mongodb.org/mongo-driver/mongo"

type AuthorController struct {
	db *mongo.Client
}

func NewAuthorController(db *mongo.Client) *AuthorController {
	return &AuthorController{
		db: db,
	}
}

type BookController struct {
	db *mongo.Client
}

func NewBookController(db *mongo.Client) *BookController {
	return &BookController{
		db: db,
	}
}
