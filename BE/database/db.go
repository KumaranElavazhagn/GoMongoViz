package database

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectMongo establishes a connection to the MongoDB Atlas cluster
// It configures the connection with the proper credentials and returns
// a connected MongoDB client that can be used throughout the application
func ConnectMongo() *mongo.Client {
	// Get credentials from environment variables
	mongoUser := "Kumarane2000"
	mongoPass := url.QueryEscape("Sathiya1981@") // Auto-encode special chars

	// Construct the MongoDB connection string with proper formatting
	mongoURI := fmt.Sprintf(
		"mongodb+srv://%s:%s@gomongoviz.6omou9u.mongodb.net/?retryWrites=true&w=majority&appName=gomongoviz",
		mongoUser, mongoPass,
	)

	// Configure client options and connect to MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Verify the connection by pinging the database
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client
}
