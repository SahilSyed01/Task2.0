// middleware_test.go
package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/ShreerajShettyK/cognitoJwtAuthenticator"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// Mocks for Secrets Manager and Cognito
// Mocks for Secrets Manager and Cognito
type mockSecretsManagerClient struct {
	MockSecretString        string
	MockGetSecretValueError bool
}

func (m *mockSecretsManagerClient) GetSecretValue(ctx context.Context, input *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	if m.MockGetSecretValueError {
		return nil, errors.New("mock error: GetSecretValue failed")
	}

	// Mock secret value
	secretString := `{"USER_POOL_ID": "mock-user-pool-id", "REGION": "mock-region"}`
	if m.MockSecretString != "" {
		secretString = m.MockSecretString
	}

	output := &secretsmanager.GetSecretValueOutput{
		SecretString: aws.String(secretString),
	}
	return output, nil
}

// Mock function for Cognito token validation
func mockCognitoValidate(ctx context.Context, region string, userPoolID string, tokenString string) (*cognitoJwtAuthenticator.AWSCognitoClaims, error) {
	// Mock token validation (always return true for simplicity)
	return &cognitoJwtAuthenticator.AWSCognitoClaims{}, nil
}

func init() {
	// Override environment variables for testing
	os.Setenv("REGION", "mock-region")
	os.Setenv("SECRET", "mock-secret-name")
}

// TestAuthenticateMiddleware tests the Authenticate middleware function
func TestAuthenticateMiddleware(t *testing.T) {
	// Set up a mock handler to simulate the next handler in the chain
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Mock initialization
	secretsClient = &mockSecretsManagerClient{}
	cognitovalidate = mockCognitoValidate

	// Create a request with a valid token
	reqValidToken, _ := http.NewRequest("GET", "/", nil)
	reqValidToken.Header.Set("Authorization", "Bearer valid-token")

	// Create a request with missing authorization header
	reqMissingHeader, _ := http.NewRequest("GET", "/", nil)

	// Create a request with invalid authorization header format
	reqInvalidHeaderFormat, _ := http.NewRequest("GET", "/", nil)
	reqInvalidHeaderFormat.Header.Set("Authorization", "invalid-format")

	// Create a test server with the Authenticate middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Authenticate(mockHandler).ServeHTTP(w, r)
	})

	// Test case: valid token
	rrValidToken := httptest.NewRecorder()
	handler.ServeHTTP(rrValidToken, reqValidToken)
	if rrValidToken.Code != http.StatusOK {
		t.Errorf("Valid token test: expected status %d, got %d", http.StatusOK, rrValidToken.Code)
	}

	// Test case: missing authorization header
	rrMissingHeader := httptest.NewRecorder()
	handler.ServeHTTP(rrMissingHeader, reqMissingHeader)
	if rrMissingHeader.Code != http.StatusUnauthorized {
		t.Errorf("Missing header test: expected status %d, got %d", http.StatusUnauthorized, rrMissingHeader.Code)
	}

	// Test case: invalid authorization header format
	rrInvalidHeaderFormat := httptest.NewRecorder()
	handler.ServeHTTP(rrInvalidHeaderFormat, reqInvalidHeaderFormat)
	if rrInvalidHeaderFormat.Code != http.StatusUnauthorized {
		t.Errorf("Invalid header format test: expected status %d, got %d", http.StatusUnauthorized, rrInvalidHeaderFormat.Code)
	}
}

// Test case for error when fetching secret value
func TestAuthenticateErrorFetchingSecret(t *testing.T) {
	// Set up a mock handler to simulate the next handler in the chain
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Mock initialization
	secretsClient = &mockSecretsManagerClient{}
	cognitovalidate = mockCognitoValidate

	// Create a request with a valid token
	reqValidToken, _ := http.NewRequest("GET", "/", nil)
	reqValidToken.Header.Set("Authorization", "Bearer valid-token")

	// Simulate an error in fetching the secret value
	mockSecretsClient := secretsClient.(*mockSecretsManagerClient)
	mockSecretsClient.MockGetSecretValueError = true

	// Create a test server with the Authenticate middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Authenticate(mockHandler).ServeHTTP(w, r)
	})

	// Perform request
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, reqValidToken)

	// Verify response status code
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
	// Verify error message in response
	expectedError := "Internal server error"
	if body := rr.Body.String(); !strings.Contains(body, expectedError) {
		t.Errorf("Expected response body to contain '%s', got '%s'", expectedError, body)
	}
}

// Test case for error when unmarshalling secret string
func TestAuthenticateErrorUnmarshalSecret(t *testing.T) {
	// Set up a mock handler to simulate the next handler in the chain
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Mock initialization
	secretsClient = &mockSecretsManagerClient{
		MockSecretString: "invalid-json", // Simulate invalid JSON in secret string
	}
	cognitovalidate = mockCognitoValidate

	// Create a request with a valid token
	reqValidToken, _ := http.NewRequest("GET", "/", nil)
	reqValidToken.Header.Set("Authorization", "Bearer valid-token")

	// Create a test server with the Authenticate middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Authenticate(mockHandler).ServeHTTP(w, r)
	})

	// Perform request
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, reqValidToken)

	// Verify response status code
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
	// Verify error message in response
	expectedError := "Failed to unmarshal secret"
	if body := rr.Body.String(); !strings.Contains(body, expectedError) {
		t.Errorf("Expected response body to contain '%s', got '%s'", expectedError, body)
	}
}

// Test case for error in token validation
func TestAuthenticateErrorTokenValidation(t *testing.T) {
	// Set up a mock handler to simulate the next handler in the chain
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Mock initialization
	secretsClient = &mockSecretsManagerClient{}
	cognitovalidate = func(ctx context.Context, region string, userPoolID string, tokenString string) (*cognitoJwtAuthenticator.AWSCognitoClaims, error) {
		return nil, errors.New("mock error: token validation failed")
	}

	// Create a request with a valid token
	reqValidToken, _ := http.NewRequest("GET", "/", nil)
	reqValidToken.Header.Set("Authorization", "Bearer invalid-token")

	// Create a test server with the Authenticate middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Authenticate(mockHandler).ServeHTTP(w, r)
	})

	// Perform request
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, reqValidToken)

	// Verify response status code
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
	// Verify error message in response
	expectedError := "Unauthorized"
	if body := rr.Body.String(); !strings.Contains(body, expectedError) {
		t.Errorf("Expected response body to contain '%s', got '%s'", expectedError, body)
	}
}

