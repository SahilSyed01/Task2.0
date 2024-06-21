package controllers

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"testing"
// 	"time"

// 	"go-chat-app/helpers"
// 	"go-chat-app/models"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// // MockUserCollection mocks MongoDB collection for testing
// type MockUserCollection struct {
// 	users map[string]models.User
// }

// // FindOne mocks the FindOne method of the MongoDB collection
// func (m *MockUserCollection) FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult {
// 	email := filter.(bson.M)["email"].(string)
// 	user, ok := m.users[email]
// 	if !ok {
// 		return mongo.NewSingleResultFromDocument(nil, errors.New("user not found"), nil)
// 	}
// 	return mongo.NewSingleResultFromDocument(user, nil, nil)
// }

// func TestLogin(t *testing.T) {
// 	// Set environment variables for testing
// 	os.Setenv("REGION", "us-west-2")
// 	os.Setenv("SECRET", "mySecret")

// 	// Mock secrets
// 	GetSecret = func(region, secretName string) (*SecretResult, error) {
// 		return &SecretResult{
// 			UserPoolID: "us-west-2_testpool",
// 			Region:     "us-west-2",
// 		}, nil
// 	}

// 	// Mock user collection
// 	mockCollection := &MockUserCollection{
// 		users: map[string]models.User{
// 			"test@example.com": {
// 				Email:    "test@example.com",
// 				Password: "hashedpassword", // Assume this is a hashed password
// 				User_id:  "12345",
// 				First_name: "Test",
// 			},
// 		},
// 	}

// 	// Mock password verification
// 	VerifyPassword = func(password, hashedPassword string) (bool, string) {
// 		if password == "password" && hashedPassword == "hashedpassword" {
// 			return true, "Password is valid"
// 		}
// 		return false, "Invalid password"
// 	}

// 	// Mock token generation
// 	helpers.GenerateToken = func(firstName, userID string) (string, error) {
// 		return "mocktoken", nil
// 	}

// 	t.Run("Successful login", func(t *testing.T) {
// 		requestBody, _ := json.Marshal(map[string]string{
// 			"email":    "test@example.com",
// 			"password": "password",
// 		})

// 		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
// 		rr := httptest.NewRecorder()

// 		// Run the Login function
// 		Login(rr, req)

// 		// Check the response status code
// 		if status := rr.Code; status != http.StatusOK {
// 			t.Errorf("handler returned wrong status code: got %v want %v",
// 				status, http.StatusOK)
// 		}

// 		// Check the response body
// 		expectedResponse := `{"Success":"True"}`
// 		if body := rr.Body.String(); body != expectedResponse {
// 			t.Errorf("handler returned unexpected body: got %v want %v",
// 				body, expectedResponse)
// 		}
// 	})

// 	t.Run("Invalid email", func(t *testing.T) {
// 		requestBody, _ := json.Marshal(map[string]string{
// 			"email":    "invalid@example.com",
// 			"password": "password",
// 		})

// 		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
// 		rr := httptest.NewRecorder()

// 		// Run the Login function
// 		Login(rr, req)

// 		// Check the response status code
// 		if status := rr.Code; status != http.StatusUnauthorized {
// 			t.Errorf("handler returned wrong status code: got %v want %v",
// 				status, http.StatusUnauthorized)
// 		}

// 		// Check the response body
// 		expectedResponse := "email or password is incorrect"
// 		if body := rr.Body.String(); body != expectedResponse {
// 			t.Errorf("handler returned unexpected body: got %v want %v",
// 				body, expectedResponse)
// 		}
// 	})

// 	t.Run("Invalid password", func(t *testing.T) {
// 		requestBody, _ := json.Marshal(map[string]string{
// 			"email":    "test@example.com",
// 			"password": "wrongpassword",
// 		})

// 		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
// 		rr := httptest.NewRecorder()

// 		// Run the Login function
// 		Login(rr, req)

// 		// Check the response status code
// 		if status := rr.Code; status != http.StatusUnauthorized {
// 			t.Errorf("handler returned wrong status code: got %v want %v",
// 				status, http.StatusUnauthorized)
// 		}

// 		// Check the response body
// 		expectedResponse := "Invalid password"
// 		if body := rr.Body.String(); body != expectedResponse {
// 			t.Errorf("handler returned unexpected body: got %v want %v",
// 				body, expectedResponse)
// 		}
// 	})

// 	t.Run("Bad request body", func(t *testing.T) {
// 		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(nil))
// 		rr := httptest.NewRecorder()

// 		// Run the Login function
// 		Login(rr, req)

// 		// Check the response status code
// 		if status := rr.Code; status != http.StatusBadRequest {
// 			t.Errorf("handler returned wrong status code: got %v want %v",
// 				status, http.StatusBadRequest)
// 		}

// 		// Check the response body
// 		expectedResponse := "EOF"
// 		if body := rr.Body.String(); body != expectedResponse {
// 			t.Errorf("handler returned unexpected body: got %v want %v",
// 				body, expectedResponse)
// 		}
// 	})
// }
