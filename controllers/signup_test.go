package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-chat-app/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	//"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mock functions
func MockHashPassword(password string) string {
	return password // Mock hash function returns the same password
}

func MockInsertOneUser(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return &mongo.InsertOneResult{}, nil // Mock successful insertion
}

func MockCountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return 0, nil // Mock count documents returns no existing documents
}

// MockGetSecret1 mocks the function to retrieve secrets from AWS Secrets Manager.
func MockGetSecret1(client SecretsManagerClient, secretName string) (*models.SecretsManagerSecret, error) {
	return &models.SecretsManagerSecret{
		UserPoolID: "test",
		Region:     "test",
	}, nil
}

func MockGetSecretFailure1(client SecretsManagerClient, secretName string) (*models.SecretsManagerSecret, error) {
	return nil, errors.New("simulated secret retrieval failure")
}

func MockGetAWSConfig1() (aws.Config, error) {
	return aws.Config{
		Region: "test",
	}, nil
}

// MockAuthenticate mocks the authentication middleware
func MockAuthenticate1(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Mock token validation or authorization logic
		// Example: Check if the request has a valid Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Example: Simulate valid token scenario
		// You can customize this based on your token validation logic
		if authHeader != "Bearer dummytoken" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Proceed with the request if token is valid
		next.ServeHTTP(w, r)
	}
}

func TestSignup_Success(t *testing.T) {
	t.Run("test for successful signup", func(t *testing.T) {
		// Create a sample user request body with required fields
		user := models.User{
			ID:         primitive.NewObjectID(),
			First_name: "John",
			Last_name:  "Doe",
			Password:   "password123",
			Email:      "test@example.com",
			Phone:      "1234567890",
			User_id:    primitive.NewObjectID().Hex(),
		}
		requestBody, _ := json.Marshal(user)

		// Create a request
		req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		req.Header = http.Header{
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer dummytoken"}, // Mocking the Authorization header with a dummy token
		}

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		// Mock dependencies
		hashPassword = MockHashPassword
		insertOneUser = MockInsertOneUser
		countDocs = MockCountDocuments
		getAWSConfig = MockGetAWSConfig1
		getSMClient = MockGetSMClient
		getSecret = MockGetSecret1
		authenticate = MockAuthenticate1 // Use the mock token validation function

		// Call the Signup handler function directly
		Signup(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Optionally, check the response body or headers if needed
	})
}

func TestSignup_MissingAuthorization(t *testing.T) {
	t.Run("test for missing authorization header", func(t *testing.T) {
		user := models.User{
			ID:         primitive.NewObjectID(),
			First_name: "John",
			Last_name:  "Doe",
			Password:   "password123",
			Email:      "test@example.com",
			Phone:      "1234567890",
			User_id:    primitive.NewObjectID().Hex(),
		}
		requestBody, _ := json.Marshal(user)

		req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		// Mock dependencies
		authenticate = MockAuthenticate1 // Use the mock token validation function

		Signup(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})
}

func TestSignup_InvalidAuthorization(t *testing.T) {
	t.Run("test for invalid authorization header", func(t *testing.T) {
		user := models.User{
			ID:         primitive.NewObjectID(),
			First_name: "John",
			Last_name:  "Doe",
			Password:   "password123",
			Email:      "test@example.com",
			Phone:      "1234567890",
			User_id:    primitive.NewObjectID().Hex(),
		}
		requestBody, _ := json.Marshal(user)

		req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		req.Header = http.Header{
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer invalidtoken"},
		}

		rr := httptest.NewRecorder()

		// Mock dependencies
		authenticate = MockAuthenticate1 // Use the mock token validation function

		Signup(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})
}
func MockCountDocumentsDuplicate(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return 1, nil // Mock count documents returns an existing document
}

// func TestSignup_DuplicateUser(t *testing.T) {
// 	t.Run("test for duplicate user", func(t *testing.T) {
// 		user := models.User{
// 			ID:         primitive.NewObjectID(),
// 			First_name: "John",
// 			Last_name:  "Doe",
// 			Password:   "password123",
// 			Email:      "test@example.com",
// 			Phone:      "1234567890",
// 			User_id:    primitive.NewObjectID().Hex(),
// 		}
// 		requestBody, _ := json.Marshal(user)

// 		req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBody))
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		req.Header = http.Header{
// 			"Content-Type":  {"application/json"},
// 			"Authorization": {"Bearer dummytoken"},
// 		}

// 		rr := httptest.NewRecorder()

// 		// Mock dependencies
// 		hashPassword = MockHashPassword
// 		insertOneUser = MockInsertOneUser
// 		countDocs = MockCountDocumentsDuplicate
// 		getAWSConfig = MockGetAWSConfig1
// 		getSMClient = MockGetSMClient
// 		getSecret = MockGetSecret1
// 		authenticate = MockAuthenticate1

// 		Signup(rr, req)

// 		if status := rr.Code; status != http.StatusConflict {
// 			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusConflict)
// 		}
// 	})
// }

func TestSignup_SecretRetrievalFailure(t *testing.T) {
	t.Run("test for secret retrieval failure", func(t *testing.T) {
		user := models.User{
			ID:         primitive.NewObjectID(),
			First_name: "John",
			Last_name:  "Doe",
			Password:   "password123",
			Email:      "test@example.com",
			Phone:      "1234567890",
			User_id:    primitive.NewObjectID().Hex(),
		}
		requestBody, _ := json.Marshal(user)

		req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		req.Header = http.Header{
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer dummytoken"},
		}

		rr := httptest.NewRecorder()

		// Mock dependencies
		hashPassword = MockHashPassword
		insertOneUser = MockInsertOneUser
		countDocs = MockCountDocuments
		getAWSConfig = MockGetAWSConfig1
		getSMClient = MockGetSMClient
		getSecret = MockGetSecretFailure1
		authenticate = MockAuthenticate1

		Signup(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}
func MockInsertOneUserFailure(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return nil, errors.New("simulated insertion failure")
}
func TestSignup_InsertionFailure(t *testing.T) {
	t.Run("test for insertion failure", func(t *testing.T) {
		user := models.User{
			ID:         primitive.NewObjectID(),
			First_name: "John",
			Last_name:  "Doe",
			Password:   "password123",
			Email:      "test@example.com",
			Phone:      "1234567890",
			User_id:    primitive.NewObjectID().Hex(),
		}
		requestBody, _ := json.Marshal(user)

		req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		req.Header = http.Header{
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer dummytoken"},
		}

		rr := httptest.NewRecorder()

		// Mock dependencies
		hashPassword = MockHashPassword
		insertOneUser = MockInsertOneUserFailure
		countDocs = MockCountDocuments
		getAWSConfig = MockGetAWSConfig1
		getSMClient = MockGetSMClient
		getSecret = MockGetSecret1
		authenticate = MockAuthenticate1

		Signup(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}

func MockGetAWSConfigFailure1() (aws.Config, error) {
	return aws.Config{}, errors.New("simulated AWS config retrieval failure")
}

func TestSignup_AWSConfigFailure1(t *testing.T) {
	t.Run("test for AWS config failure", func(t *testing.T) {
		// Create a sample user request body
		user := models.User{
			ID:         primitive.NewObjectID(),
			First_name: "John",
			Last_name:  "Doe",
			Password:   "password123",
			Email:      "test@example.com",
			Phone:      "1234567890",
			User_id:    primitive.NewObjectID().Hex(),
		}
		requestBody, _ := json.Marshal(user)

		// Create a request
		req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		req.Header = http.Header{
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer dummytoken"}, // Mocking the Authorization header with a dummy token
		}

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		// Mock dependencies
		hashPassword = MockHashPassword
		insertOneUser = MockInsertOneUser
		countDocs = MockCountDocuments
		getAWSConfig = MockGetAWSConfigFailure
		getSMClient = MockGetSMClient
		getSecret = MockGetSecret1
		authenticate = MockAuthenticate1 // Use the mock token validation function

		// Call the Signup handler function directly
		Signup(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}

func TestSignup_InvalidJSON(t *testing.T) {
	t.Run("test for invalid JSON", func(t *testing.T) {
		// Create an invalid JSON request body
		invalidJSON := `{"email": "test@example.com", "password": "password123", "phone": "1234567890"` // Missing closing bracket

		// Create a request
		req, err := http.NewRequest("POST", "/signup", bytes.NewBufferString(invalidJSON))
		if err != nil {
			t.Fatal(err)
		}
		req.Header = http.Header{
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer dummytoken"}, // Mocking the Authorization header with a dummy token
		}

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		// Mock dependencies
		hashPassword = MockHashPassword
		insertOneUser = MockInsertOneUser
		countDocs = MockCountDocuments
		getAWSConfig = MockGetAWSConfig1
		getSMClient = MockGetSMClient
		getSecret = MockGetSecret1
		authenticate = MockAuthenticate1 // Use the mock token validation function

		// Call the Signup handler function directly
		Signup(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}

func TestSignup_ValidationErrors(t *testing.T) {
	t.Run("test for validation errors", func(t *testing.T) {
		// Create a user request body with invalid fields
		user := models.User{
			ID:         primitive.NewObjectID(),
			First_name: "", // Missing first name
			Last_name:  "Doe",
			Password:   "123",           // Password too short
			Email:      "invalid-email", // Invalid email format
			Phone:      "",              // Missing phone number
			User_id:    primitive.NewObjectID().Hex(),
		}
		requestBody, _ := json.Marshal(user)

		// Create a request
		req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		req.Header = http.Header{
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer dummytoken"}, // Mocking the Authorization header with a dummy token
		}

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		// Mock dependencies
		hashPassword = MockHashPassword
		insertOneUser = MockInsertOneUser
		countDocs = MockCountDocuments
		getAWSConfig = MockGetAWSConfig1
		getSMClient = MockGetSMClient
		getSecret = MockGetSecret1
		authenticate = MockAuthenticate1 // Use the mock token validation function

		// Call the Signup handler function directly
		Signup(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}

func MockCountDocumentsFailure(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return 0, errors.New("simulated count documents failure")
}

func TestSignup_CountDocumentsEmailFailure(t *testing.T) {
	t.Run("test for count documents email failure", func(t *testing.T) {
		// Create a sample user request body
		user := models.User{
			ID:         primitive.NewObjectID(),
			First_name: "John",
			Last_name:  "Doe",
			Password:   "password123",
			Email:      "test@example.com",
			Phone:      "1234567890",
			User_id:    primitive.NewObjectID().Hex(),
		}
		requestBody, _ := json.Marshal(user)

		// Create a request
		req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		req.Header = http.Header{
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer dummytoken"}, // Mocking the Authorization header with a dummy token
		}

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		// Mock dependencies
		hashPassword = MockHashPassword
		insertOneUser = MockInsertOneUser
		countDocs = MockCountDocumentsFailure // Simulate failure in counting documents by email
		getAWSConfig = MockGetAWSConfig1
		getSMClient = MockGetSMClient
		getSecret = MockGetSecret1
		authenticate = MockAuthenticate1 // Use the mock token validation function

		// Call the Signup handler function directly
		Signup(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}


func TestSignup_SecretIncomplete(t *testing.T) {
    t.Run("test for secret incomplete", func(t *testing.T) {
        // Create a sample user request body with required fields
        user := models.User{
            ID:         primitive.NewObjectID(),
            First_name: "John",
            Last_name:  "Doe",
            Password:   "password123",
            Email:      "test@example.com",
            Phone:      "1234567890",
            User_id:    primitive.NewObjectID().Hex(),
        }
        requestBody, _ := json.Marshal(user)

        // Create a request
        req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBody))
        if err != nil {
            t.Fatal(err)
        }
        req.Header = http.Header{
            "Content-Type":  {"application/json"},
            "Authorization": {"Bearer dummytoken"}, // Mocking the Authorization header with a dummy token
        }

        // Mock dependencies
        hashPassword = MockHashPassword
        insertOneUser = MockInsertOneUser
        countDocs = MockCountDocuments
        getAWSConfig = MockGetAWSConfig1
        getSMClient = MockGetSMClient
        getSecret = func(client SecretsManagerClient, secretName string) (*models.SecretsManagerSecret, error) {
            // Simulate fetching a secret that is incomplete or nil
            return nil, nil // Mock returning nil secret
        }
        authenticate = MockAuthenticate1 // Use the mock token validation function

        defer func() {
            // Restore original functions after the test
            getSecret = MockGetSecret1
            authenticate = MockAuthenticate1
        }()

        // Create a ResponseRecorder to record the response
        rr := httptest.NewRecorder()

        // Call the Signup handler function directly
        Signup(rr, req)

        // Check the status code
        if status := rr.Code; status != http.StatusInternalServerError {
            t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
        }

        // Optionally, check the response body or logs if needed
        // For example, you could check the response body contains "Internal server error"
        responseBody, err := ioutil.ReadAll(rr.Body)
        if err != nil {
            t.Fatal(err)
        }
        expectedErrorMessage := "Internal server error"
        if !strings.Contains(string(responseBody), expectedErrorMessage) {
            t.Errorf("expected error message '%s' not found in response body: %s", expectedErrorMessage, responseBody)
        }
    })
}

func TestSignup_CountDocumentsPhoneError(t *testing.T) {
    t.Run("test for error counting documents for phone number", func(t *testing.T) {
        // Create a sample user request body with required fields
        user := models.User{
            ID:         primitive.NewObjectID(),
            First_name: "John",
            Last_name:  "Doe",
            Password:   "password123",
            Email:      "test@example.com",
            Phone:      "123456890",
            User_id:    primitive.NewObjectID().Hex(),
        }
        requestBody, _ := json.Marshal(user)

        // Create a request
        req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBody))
        if err != nil {
            t.Fatal(err)
        }
        req.Header = http.Header{
            "Content-Type":  {"application/json"},
            "Authorization": {"Bearer dummytoken"}, // Mocking the Authorization header with a dummy token
        }

        // Create a ResponseRecorder to record the response
        rr := httptest.NewRecorder()

        // Mock dependencies
        hashPassword = MockHashPassword
        insertOneUser = MockInsertOneUser
        countDocs = func(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
            return 0, errors.New("simulated count documents failure")
        }
        getAWSConfig = MockGetAWSConfig1
        getSMClient = MockGetSMClient
        getSecret = MockGetSecret1
        authenticate = MockAuthenticate1 // Use the mock token validation function

        // Call the Signup handler function directly
        Signup(rr, req)

        // Check the status code
        if status := rr.Code; status != http.StatusInternalServerError {
            t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
        }

        // Check the response body for the expected error message
        expectedError := "error occurred while checking for the phone number"
        if !strings.Contains(rr.Body.String(), expectedError) {
            t.Errorf("expected error message '%s' not found in response body: %s", expectedError, rr.Body.String())
        }
    })
}


func TestSignup_PhoneNumberAlreadyExists(t *testing.T) {
    t.Run("test for existing phone number", func(t *testing.T) {
        // Create a sample user request body with required fields
        user := models.User{
            ID:         primitive.NewObjectID(),
            First_name: "John",
            Last_name:  "Doe",
            Password:   "password123",
            Email:      "test@example.com",
            Phone:      "1234567890",
            User_id:    primitive.NewObjectID().Hex(),
        }
        requestBody, _ := json.Marshal(user)

        // Create a request
        req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBody))
        if err != nil {
            t.Fatal(err)
        }
        req.Header = http.Header{
            "Content-Type":  {"application/json"},
            "Authorization": {"Bearer dummytoken"}, // Mocking the Authorization header with a dummy token
        }

        // Create a ResponseRecorder to record the response
        rr := httptest.NewRecorder()

        // Mock dependencies
        hashPassword = MockHashPassword
        insertOneUser = MockInsertOneUser
        countDocs = func(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
            return 1, nil // Simulate that the phone number already exists
        }
        getAWSConfig = MockGetAWSConfig1
        getSMClient = MockGetSMClient
        getSecret = MockGetSecret1
        authenticate = MockAuthenticate1 // Use the mock token validation function

        // Call the Signup handler function directly
        Signup(rr, req)

        // Check the status code
        if status := rr.Code; status != http.StatusBadRequest {
            t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
        }

        // Check the response body for the expected error message
        expectedError := "this phone number already exists"
        if !strings.Contains(rr.Body.String(), expectedError) {
            t.Errorf("expected error message '%s' not found in response body: %s", expectedError, rr.Body.String())
        }
    })
}

