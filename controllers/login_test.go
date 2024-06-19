package controllers

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"go-chat-app/models"

// 	"github.com/stretchr/testify/assert"
// 	// "go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// // MockCollection implementation for testing
// type MockCollection struct {
// 	mockData map[string]interface{}
// }

// func (c *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
// 	// Implement mock FindOne logic here
// 	// For simplicity, returning nil for now
// 	return nil
// }

// func (c *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
// 	// Implement mock InsertOne logic here
// 	// For simplicity, returning nil for now
// 	return nil, nil
// }

// func TestLogin(t *testing.T) {
// 	// Mock setup
// 	mockUser := &models.User{
// 		Email:    stringPtr("test@example.com"),
// 		Password: stringPtr("mockPassword"),
// 	}

// 	// Mock MongoDB collection
// 	mockUserCollection := new(MockCollection)

// 	// Replace actual implementation with mock
// 	userCollection = mockUserCollection

// 	// Create a request
// 	requestBody, _ := json.Marshal(mockUser)
// 	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Recorder for response
// 	rr := httptest.NewRecorder()

// 	// Call the Login function directly (since it's a handler function)
// 	Login(rr, req)

// 	// Assert the status code
// 	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200")

// 	// Assert the response body
// 	expectedBody := `{"Success":"True"}`
// 	assert.Equal(t, expectedBody, rr.Body.String(), "Response body should match expected")
// }

// // Utility function for creating string pointers
// func stringPtr(s string) *string {
// 	return &s
// }
