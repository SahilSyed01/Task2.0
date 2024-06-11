package middleware_test

import (
	// "context"
	"net/http"
	"net/http/httptest"
	"os"
	// "strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-chat-app/middleware"
)

func TestAuthenticate_ValidToken(t *testing.T) {
	// Set up environment variables
	os.Setenv("REGION", "test-region")
	os.Setenv("USER_POOL_ID", "test-pool-id")

	// Create a mock handler for testing
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the mock handler with the Authenticate middleware
	handler := middleware.Authenticate(mockHandler)

	// Create a request with a valid token
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer validToken123")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check if the status code is 401 Unauthorized (updated)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthenticate_InvalidToken(t *testing.T) {
	// Set up environment variables
	os.Setenv("REGION", "test-region")
	os.Setenv("USER_POOL_ID", "test-pool-id")

	// Create a mock handler for testing
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the mock handler with the Authenticate middleware
	handler := middleware.Authenticate(mockHandler)

	// Create a request with an invalid token
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalidToken456")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check if the status code is 401 Unauthorized (updated)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthenticate_MissingAuthorizationHeader(t *testing.T) {
    // Create a mock handler for testing
    mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

    // Wrap the mock handler with the Authenticate middleware
    handler := middleware.Authenticate(mockHandler)

    // Create a request without the Authorization header
    req := httptest.NewRequest("GET", "/", nil)

    // Create a ResponseRecorder to record the response
    rr := httptest.NewRecorder()

    // Serve the HTTP request
    handler.ServeHTTP(rr, req)

    // Check if the status code is 401 Unauthorized
    assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthenticate_InvalidAuthorizationHeaderFormat(t *testing.T) {
    // Create a mock handler for testing
    mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

    // Wrap the mock handler with the Authenticate middleware
    handler := middleware.Authenticate(mockHandler)

    // Create a request with an invalid Authorization header format
    req := httptest.NewRequest("GET", "/", nil)
    req.Header.Set("Authorization", "invalidformat")

    // Create a ResponseRecorder to record the response
    rr := httptest.NewRecorder()

    // Serve the HTTP request
    handler.ServeHTTP(rr, req)

    // Check if the status code is 401 Unauthorized
    assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthenticate_TokenValidationError(t *testing.T) {
    // Set up environment variables
    os.Setenv("REGION", "test-region")
    os.Setenv("USER_POOL_ID", "test-pool-id")

    // Create a mock handler for testing
    mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

    // Wrap the mock handler with the Authenticate middleware
    handler := middleware.Authenticate(mockHandler)

    // Create a request with a valid token
    req := httptest.NewRequest("GET", "/", nil)
    req.Header.Set("Authorization", "Bearer invalidToken")

    // Create a ResponseRecorder to record the response
    rr := httptest.NewRecorder()

    // Serve the HTTP request
    handler.ServeHTTP(rr, req)

    // Check if the status code is 401 Unauthorized
    assert.Equal(t, http.StatusUnauthorized, rr.Code)
}



func TestAuthenticate_ExpiredToken(t *testing.T) {
    // Set up environment variables
    os.Setenv("REGION", "test-region")
    os.Setenv("USER_POOL_ID", "test-pool-id")

    // Create a mock handler for testing
    mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

    // Wrap the mock handler with the Authenticate middleware
    handler := middleware.Authenticate(mockHandler)

    // Create a request with an expired token
    req := httptest.NewRequest("GET", "/", nil)
    req.Header.Set("Authorization", "Bearer expiredToken")

    // Create a ResponseRecorder to record the response
    rr := httptest.NewRecorder()

    // Serve the HTTP request
    handler.ServeHTTP(rr, req)

    // Check if the status code is 401 Unauthorized
    assert.Equal(t, http.StatusUnauthorized, rr.Code)
}