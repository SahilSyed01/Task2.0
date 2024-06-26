package database

import (
	"context"
	"encoding/json"
	"errors"
	//"log"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	//"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)
type MockCollection struct {
    // Add fields as needed for mock implementation
}

func (c *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
    // Implement mock FindOne logic
    return nil
}

func (c *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
    // Implement mock InsertOne logic
    return nil, nil
}

// NewMockCollection creates a new instance of MockCollection
func NewMockCollection() *MockCollection {
    return &MockCollection{}
}
// AWSClient is an interface for AWS Secrets Manager client.

// MockAWSClient is a mock implementation of the AWS Secrets Manager client.
type MockAWSClient struct{
	GetSecretValueFunc func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

func TestDBinstance(t *testing.T) {
	// Backup original functions
	origSecretsManagerClient := secretsManagerClient
	origCreateMongoClient := createMongoClient
	origConnectMongoClient := connectMongoClient
	origPingMongoClient := pingMongoClient

	defer func() {
		// Restore original functions
		secretsManagerClient = origSecretsManagerClient
		createMongoClient = origCreateMongoClient
		connectMongoClient = origConnectMongoClient
		pingMongoClient = origPingMongoClient
	}()

	t.Run("Simulated error retrieving secret", func(t *testing.T) {
		// Set up a mock AWS Secrets Manager client
		secretsManagerClient = &MockAWSClient{}

		// Set simulateError to true
		simulateError = true
		client := DBinstance()
		assert.Nil(t, client, "Expected DBinstance to return nil when simulateError is true")
	})

	t.Run("Invalid JSON in secret", func(t *testing.T) {
		// Set up a mock AWS Secrets Manager client
		secretsManagerClient = &MockAWSClient{
			GetSecretValueFunc: func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
				return &secretsmanager.GetSecretValueOutput{
					SecretString: aws.String("invalid JSON"),
				}, nil
			},
		}

		// Set simulateError to false
		simulateError = false
		client := DBinstance()
		assert.Nil(t, client, "Expected DBinstance to return nil with invalid JSON in secret")
	})

	t.Run("Error creating MongoDB client", func(t *testing.T) {
		// Set up a mock AWS Secrets Manager client
		secretsManagerClient = &MockAWSClient{
			GetSecretValueFunc: func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
				secret := map[string]string{"connectionString": "mongodb://localhost:27017"}
				secretBytes, _ := json.Marshal(secret)
				return &secretsmanager.GetSecretValueOutput{
					SecretString: aws.String(string(secretBytes)),
				}, nil
			},
		}

		// Mock createMongoClient to return error
		createMongoClient = func(opts ...*options.ClientOptions) (*mongo.Client, error) {
			return nil, errors.New("error creating MongoDB client")
		}

		client := DBinstance()
		assert.Nil(t, client, "Expected DBinstance to return nil on error creating MongoDB client")
	})

	t.Run("Error connecting to MongoDB", func(t *testing.T) {
		// Set up a mock AWS Secrets Manager client
		secretsManagerClient = &MockAWSClient{
			GetSecretValueFunc: func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
				secret := map[string]string{"connectionString": "mongodb://localhost:27017"}
				secretBytes, _ := json.Marshal(secret)
				return &secretsmanager.GetSecretValueOutput{
					SecretString: aws.String(string(secretBytes)),
				}, nil
			},
		}

		// Mock createMongoClient and connectMongoClient to return no error
		createMongoClient = func(opts ...*options.ClientOptions) (*mongo.Client, error) {
			return &mongo.Client{}, nil
		}
		connectMongoClient = func(client *mongo.Client, ctx context.Context) error {
			return errors.New("error connecting to MongoDB")
		}

		client := DBinstance()
		assert.Nil(t, client, "Expected DBinstance to return nil on error connecting to MongoDB")
	})

	t.Run("Error pinging MongoDB server", func(t *testing.T) {
		// Set up a mock AWS Secrets Manager client
		secretsManagerClient = &MockAWSClient{
			GetSecretValueFunc: func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
				secret := map[string]string{"connectionString": "mongodb://localhost:27017"}
				secretBytes, _ := json.Marshal(secret)
				return &secretsmanager.GetSecretValueOutput{
					SecretString: aws.String(string(secretBytes)),
				}, nil
			},
		}

		// Mock createMongoClient, connectMongoClient, and pingMongoClient to return no error
		createMongoClient = func(opts ...*options.ClientOptions) (*mongo.Client, error) {
			return &mongo.Client{}, nil
		}
		connectMongoClient = func(client *mongo.Client, ctx context.Context) error {
			return nil
		}
		pingMongoClient = func(client *mongo.Client, ctx context.Context, rp *readpref.ReadPref) error {
			return errors.New("error pinging MongoDB server")
		}

		client := DBinstance()
		assert.Nil(t, client, "Expected DBinstance to return nil on error pinging MongoDB server")
	})

	t.Run("Successful connection", func(t *testing.T) {
		// Set up a valid JSON for testing
		secretsManagerClient = &MockAWSClient{
			GetSecretValueFunc: func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
				secret := map[string]string{"connectionString": "mongodb://localhost:27017"}
				secretBytes, _ := json.Marshal(secret)
				return &secretsmanager.GetSecretValueOutput{
					SecretString: aws.String(string(secretBytes)),
				}, nil
			},
		}

		// Mock MongoDB client creation and connection
		createMongoClient = func(opts ...*options.ClientOptions) (*mongo.Client, error) {
			return &mongo.Client{}, nil
		}
		connectMongoClient = func(client *mongo.Client, ctx context.Context) error {
			return nil
		}
		pingMongoClient = func(client *mongo.Client, ctx context.Context, rp *readpref.ReadPref) error {
			return nil
		}

		client := DBinstance()
		assert.NotNil(t, client, "Expected DBinstance to return a non-nil client on successful connection")
	})
	// t.Run("Error loading AWS config", func(t *testing.T) {
	// 	// Mock AWS config loading to return an error
	// 	secretsManagerClient = nil // Force re-initialization in DBinstance
	// 	simulateError = false     // Ensure no simulated error for this case

	// 	// Mock AWS config loading to return an error
	// 	cfg, err := config.LoadDefaultConfig(context.Background())
	// 	if err != nil {
	// 		log.Println("Error loading AWS config:", err)
	// 		assert.Nil(t, cfg, "Expected AWS config to be nil on error")
	// 		secretsManagerClient = nil // Set to nil to force re-initialization
	// 		return
	// 	}
	// 	secretsManagerClient = secretsmanager.NewFromConfig(cfg)

	// 	client := DBinstance()
	// 	assert.Nil(t, client, "Expected DBinstance to return nil on error loading AWS config")
	// })
	t.Run("Error retrieving secret", func(t *testing.T) {
		// Set up a mock AWS Secrets Manager client
		secretsManagerClient = &MockAWSClient{
			GetSecretValueFunc: func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
				return nil, errors.New("simulated error retrieving secret")
			},
		}

		client := DBinstance()
		assert.Nil(t, client, "Expected DBinstance to return nil on error retrieving secret")
	})

	t.Run("Connection string not found in secret", func(t *testing.T) {
		// Set up a mock AWS Secrets Manager client
		secretsManagerClient = &MockAWSClient{
			GetSecretValueFunc: func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
				secret := map[string]string{}
				secretBytes, _ := json.Marshal(secret)
				return &secretsmanager.GetSecretValueOutput{
					SecretString: aws.String(string(secretBytes)),
				}, nil
			},
		}

		client := DBinstance()
		assert.Nil(t, client, "Expected DBinstance to return nil when connection string not found in secret")
	})
}


func TestOpenCollection(t *testing.T) {
	mockClient := &mongo.Client{}
	collection := OpenCollection(mockClient, "test-collection")
	assert.NotNil(t, collection, "Expected OpenCollection to return a non-nil collection")
	assert.Equal(t, "test-collection", collection.Name(), "Expected collection name to match")
}


func (m *MockAWSClient) GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	if m.GetSecretValueFunc != nil {
		return m.GetSecretValueFunc(ctx, params, optFns...)
	}
	if simulateError {
		return nil, errors.New("simulated error retrieving secret")
	}
	return &secretsmanager.GetSecretValueOutput{
		SecretString: aws.String("invalid JSON"),
	}, nil
}
