package controllers

import (
	"context"
	"encoding/json"
	"errors"

	//"errors"
	"go-chat-app/models"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	//"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// MockAuthenticate2 remains the same
func MockAuthenticate2(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}
		expectedToken := "Bearer valid-token"
		if authHeader != expectedToken {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// MockGetSecret2 remains the same
func MockGetSecret2(client SecretsManagerClient, secretName string) (*models.SecretsManagerSecret, error) {
	return &models.SecretsManagerSecret{
		UserPoolID:   "test",
		ClientID:     "test",
		ClientSecret: "test",
		Username:     "test",
		Password:     "test",
		Region:       "test",
	}, nil
}

// New mock function for MongoDB aggregation success case
func MockAggregateSuccess(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	documents := []bson.D{
		{
			{"total_count", 1},
			{"user_items", bson.A{
				bson.D{
					{"email", "test@example.com"},
					{"first_name", "Test"},
					{"last_name", "User"},
					{"password", "password"},
					{"phone", "1234567890"},
					{"user_id", "1"},
				},
			}},
		},
	}
	var interfaces []interface{}
	for _, doc := range documents {
		interfaces = append(interfaces, doc)
	}
	cursor, err := mongo.NewCursorFromDocuments(interfaces, nil, nil)
	return cursor, err
}

// New mock function for MongoDB aggregation with no users found
func MockAggregateNoUsers(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	documents := []bson.D{}
	var interfaces []interface{}
	for _, doc := range documents {
		interfaces = append(interfaces, doc)
	}
	cursor, err := mongo.NewCursorFromDocuments(interfaces, nil, nil)
	return cursor, err
}

// Test for successful response with users found
func TestGetUsers_Success(t *testing.T) {
	os.Setenv("SECRET", "test_secret")
	defer os.Unsetenv("SECRET")

	authenticate = MockAuthenticate2
	getSecret = MockGetSecret2
	aggregate = MockAggregateSuccess

	req, err := http.NewRequest("GET", "/users?recordPerPage=10&page=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer valid-token")

	rr := httptest.NewRecorder()
	GetUsers(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, expectedContentType)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("error decoding response: %v", err)
	}
	if response["total_count"].(float64) != 1 {
		t.Errorf("expected total_count 1, got %v", response["total_count"])
	}
}

// Test for no users found
func TestGetUsers_NoUsersFound(t *testing.T) {
	os.Setenv("SECRET", "test_secret")
	defer os.Unsetenv("SECRET")

	authenticate = MockAuthenticate2
	getSecret = MockGetSecret2
	aggregate = MockAggregateNoUsers

	req, err := http.NewRequest("GET", "/users?recordPerPage=10&page=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer valid-token")

	rr := httptest.NewRecorder()
	GetUsers(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func MockGetSecretWithError(client SecretsManagerClient, secretName string) (*models.SecretsManagerSecret, error) {
	return nil, errors.New("mock error fetching secret")
}

func TestGetUsers_ErrorFetchingSecret(t *testing.T) {
	os.Setenv("SECRET", "test_secret")
	defer os.Unsetenv("SECRET")

	authenticate = MockAuthenticate2
	getSecret = MockGetSecretWithError

	req, err := http.NewRequest("GET", "/users?recordPerPage=10&page=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer valid-token")

	rr := httptest.NewRecorder()
	GetUsers(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func MockGetSecretWithMissingFields(client SecretsManagerClient, secretName string) (*models.SecretsManagerSecret, error) {
	return &models.SecretsManagerSecret{
		// You can intentionally leave out UserPoolID or Region to simulate missing fields.
		ClientID:     "test",
		ClientSecret: "test",
		Username:     "test",
		Password:     "test",
	}, nil
}

func TestGetUsers_SecretMissingFields(t *testing.T) {
	os.Setenv("SECRET", "test_secret")
	defer os.Unsetenv("SECRET")

	authenticate = MockAuthenticate2
	getSecret = MockGetSecretWithMissingFields

	req, err := http.NewRequest("GET", "/users?recordPerPage=10&page=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer valid-token")

	rr := httptest.NewRecorder()
	GetUsers(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func MockAggregateError(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	return nil, errors.New("mock MongoDB aggregation error")
}

func TestGetUsers_ErrorInAggregation(t *testing.T) {
	os.Setenv("SECRET", "test_secret")
	defer os.Unsetenv("SECRET")

	authenticate = MockAuthenticate2
	getSecret = MockGetSecret2
	aggregate = MockAggregateError

	req, err := http.NewRequest("GET", "/users?recordPerPage=10&page=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer valid-token")

	rr := httptest.NewRecorder()
	GetUsers(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func MockAggregateDecodeError(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	// Simulate a MongoDB cursor with incorrect structure or data causing decoding error
	documents := []bson.D{
		{
			{"total_count", 1},
			// Simulate a structure that does not match expected decoding format
			{"user_items", bson.D{}}, // Use bson.D instead of bson.A
		},
	}
	var interfaces []interface{}
	for _, doc := range documents {
		interfaces = append(interfaces, doc)
	}
	cursor, err := mongo.NewCursorFromDocuments(interfaces, nil, nil)
	return cursor, err
}

func TestGetUsers_ErrorDecodingCursor(t *testing.T) {
	os.Setenv("SECRET", "test_secret")
	defer os.Unsetenv("SECRET")

	authenticate = MockAuthenticate2
	getSecret = MockGetSecret2
	aggregate = MockAggregateDecodeError

	req, err := http.NewRequest("GET", "/users?recordPerPage=10&page=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer valid-token")

	rr := httptest.NewRecorder()
	GetUsers(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestGetUsers_ErrorEncodingResponse(t *testing.T) {
	os.Setenv("SECRET", "test_secret")
	defer os.Unsetenv("SECRET")

	authenticate = MockAuthenticate2
	getSecret = MockGetSecret2
	aggregate = MockAggregateSuccess

	req, err := http.NewRequest("GET", "/users?recordPerPage=10&page=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer valid-token")

	// Use the custom errorResponseWriter
	rr := &errorResponseWriter{}

	GetUsers(rr, req)

	// Check the expected status code
	if status := rr.status; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

// Custom writer that always returns an error on Write
type errorResponseWriter struct {
    status int
}

func (w *errorResponseWriter) Header() http.Header {
    // Implement Header as needed (not necessary for error testing)
    return http.Header{}
}

func (w *errorResponseWriter) Write([]byte) (int, error) {
    return 0, errors.New("mock error writing")
}

func (w *errorResponseWriter) WriteHeader(statusCode int) {
    w.status = statusCode
}
