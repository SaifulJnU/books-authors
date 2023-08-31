
//for normal create, update, delete, watching code:
package main

import (
"context"
"fmt"
"go.mongodb.org/mongo-driver/bson"
"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/mongo/options"
"log"
"time"
)

func main() {
// MongoDB connection URI
uri := "mongodb://admin:secret@localhost:27017,localhost:27018,localhost:27019/?replicaSet=rs0"

	// Create a MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	fmt.Println("Connected to MongoDB!")

	// Access a database and collection
	database := client.Database("mydb")
	collection := database.Collection("mycollection")

	// Set up change stream pipeline
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"operationType": bson.M{"$in": []string{"insert", "update", "delete"}}}}},
	}
	changeStreamOpts := options.ChangeStream().SetFullDocument(options.UpdateLookup)

	// Watch for changes
	changeStream, err := collection.Watch(ctx, pipeline, changeStreamOpts)
	if err != nil {
		log.Fatal(err)
	}
	defer changeStream.Close(ctx)

	// Process change events
	for changeStream.Next(ctx) {
		var changeEvent map[string]interface{}
		if err := changeStream.Decode(&changeEvent); err != nil {
			log.Println("Error decoding change event:", err)
			continue
		}

		fmt.Println("Change event:", changeEvent)
	}
}
