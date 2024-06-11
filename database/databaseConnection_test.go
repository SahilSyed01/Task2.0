package database

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

func TestDBinstance(t *testing.T) {
	// Call the DBinstance function to get a MongoDB client
	client := DBinstance()
	defer client.Disconnect(context.Background())

	// Check if the client is nil
	if client == nil {
		t.Errorf("DBinstance() returned nil client")
	}

	// Test a simple database operation to ensure connection is successful
	err := client.Ping(context.Background(), nil)
	if err != nil {
		t.Errorf("Error pinging MongoDB: %v", err)
	}
}

func TestOpenCollection(t *testing.T) {
	// Connect to MongoDB for testing
	client := DBinstance()
	defer client.Disconnect(context.Background())

	// Test collection name
	collectionName := "testCollection"

	// Open the collection
	collection := OpenCollection(client, collectionName)

	// Check if the collection is nil
	if collection == nil {
		t.Errorf("OpenCollection() returned nil for collection %s", collectionName)
	}

	// Optional: You can further test operations on the collection if needed
}

func TestIntegration(t *testing.T) {
	// This test can be used for integration testing, making sure that DBinstance and OpenCollection work together
	client := DBinstance()
	defer client.Disconnect(context.Background())

	collectionName := "testIntegrationCollection"

	// Insert a document into the collection
	collection := OpenCollection(client, collectionName)
	if collection == nil {
		t.Fatalf("Failed to open collection %s", collectionName)
	}

	_, err := collection.InsertOne(context.Background(), bson.M{"key": "value"})
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	}
}
