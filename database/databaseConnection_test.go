package database

import (
	"context"
	// "errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MockMongoClient struct {
	*mongo.Client
	PingFunc     func(ctx context.Context, rp *readpref.ReadPref) error
	DatabaseFunc func(name string) *mongo.Database
}

func (m *MockMongoClient) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	if m.PingFunc != nil {
		return m.PingFunc(ctx, rp)
	}
	return nil
}

func (m *MockMongoClient) Database(name string) *mongo.Database {
	if m.DatabaseFunc != nil {
		return m.DatabaseFunc(name)
	}
	return nil
}

func TestDBinstance(t *testing.T) {
	tests := []struct {
		name       string
		simulate   bool
		expected   bool
		assertFunc func(*testing.T, *mongo.Client)
	}{
		{
			name:     "Success",
			simulate: false,
			expected: true,
			assertFunc: func(t *testing.T, client *mongo.Client) {
				assert.NotNil(t, client)
			},
		},
		{
			name:     "AWSConfigError",
			simulate: true,
			expected: false,
			assertFunc: func(t *testing.T, client *mongo.Client) {
				assert.Nil(t, client)
			},
		},
		// Add more test cases for other error scenarios
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simulateError = test.simulate
			client := DBinstance()
			test.assertFunc(t, client)
		})
	}
}

func TestOpenCollection(t *testing.T) {
	mockMongoClient := &MockMongoClient{
		Client: &mongo.Client{},
	}

	collection := OpenCollection(mockMongoClient.Client, "test_collection")

	assert.NotNil(t, collection)
}
