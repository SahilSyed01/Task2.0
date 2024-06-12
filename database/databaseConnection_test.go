package database

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

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

	// Set up a mock AWS Secrets Manager client
	secretsManagerClient = &MockAWSClient{}

	// Test case: Simulated error retrieving secret
	simulateError = true
	client := DBinstance()
	assert.Nil(t, client, "Expected DBinstance to return nil when simulateError is true")

	// Test case: Successful retrieval with invalid JSON
	simulateError = false
	client = DBinstance()
	assert.Nil(t, client, "Expected DBinstance to return nil with invalid JSON in secret")

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

	// Test case: Successful connection
	client = DBinstance()
	assert.NotNil(t, client, "Expected DBinstance to return a non-nil client on successful connection")
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
