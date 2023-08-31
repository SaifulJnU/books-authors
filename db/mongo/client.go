package mongo

import (
	"context"

	"github.com/saifujnu/books-authors/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// type Mongo struct {
// 	*mongo.Client
// 	database *mongo.Database
// 	name     string
// }

// var (
// 	client   *mongo.Client
// 	database *mongo.Database
// 	err      error
// )

// const (
// 	connectTimeout = 20
// )

func Connect() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(config.MongoURL)
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
