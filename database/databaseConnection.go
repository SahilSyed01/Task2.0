package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	// "testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	// "github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	secretsManagerClient AWSClient  // AWS Secrets Manager client
	simulateError        bool       // Flag to simulate error
)

// AWSClient is an interface for AWS Secrets Manager client.
type AWSClient interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

// MockAWSClient is a mock implementation of the AWS Secrets Manager client.
type MockAWSClient struct{}

func (m *MockAWSClient) GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	if simulateError {
		return nil, fmt.Errorf("simulated error retrieving secret")
	}
	// Simulate successful retrieval of secret value
	return &secretsmanager.GetSecretValueOutput{
		SecretString: aws.String("invalid JSON"),
	}, nil
}

// DBinstance connects to MongoDB using a connection string from AWS Secrets Manager.
func DBinstance() *mongo.Client {
	if secretsManagerClient == nil {
		// If secretsManagerClient is not initialized, initialize it
		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			log.Println("Error loading AWS config:", err)
			return nil
		}
		secretsManagerClient = secretsmanager.NewFromConfig(cfg)
	}

	if simulateError {
		log.Println("Simulated error loading AWS config")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Retrieve the MongoDB connection string from Secrets Manager
	secretValue, err := secretsManagerClient.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String("myApp/mongo-db-credentials"), // Replace with your secret ID
	})
	if err != nil {
		log.Println("Error retrieving secret:", err)
		return nil
	}

	// Parse the secret string to extract the connection string
	var secretsMap map[string]string
	if err := json.Unmarshal([]byte(*secretValue.SecretString), &secretsMap); err != nil {
		log.Println("Error unmarshalling secret:", err)
		return nil
	}

	connectionString, exists := secretsMap["connectionString"]
	if !exists {
		log.Println("Connection string not found in secret")
		return nil
	}

	// Create a new MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Println("Error creating MongoDB client:", err)
		return nil
	}

	// Connect to MongoDB
	err = client.Connect(ctx)
	if err != nil {
		log.Println("Error connecting to MongoDB:", err)
		return nil
	}

	// Ping the MongoDB server to ensure connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println("Error pinging MongoDB server:", err)
		return nil
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

var Client *mongo.Client

func init() {
	Client = DBinstance()
}

// OpenCollection opens a specific MongoDB collection.
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("cluster0").Collection(collectionName)
	return collection
}
