package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"go-chat-app/models"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"

	// "github.com/aws/smithy-go"

	// "github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/stretchr/testify/assert"
)
func TestAuthenticateMiddleware(t *testing.T) {
	// Mock handler to simulate the next handler in the chain
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock handler does nothing for now
		w.WriteHeader(http.StatusOK)
	})

	// Mock request with a valid JWT token
	reqValidToken := httptest.NewRequest("GET", "/test", nil)
	reqValidToken.Header.Set("Authorization", "Bearer valid_token")

	// Mock request with missing Authorization header
	reqMissingAuthHeader := httptest.NewRequest("GET", "/test", nil)

	// Mock request with invalid Authorization header format
	reqInvalidAuthHeader := httptest.NewRequest("GET", "/test", nil)
	reqInvalidAuthHeader.Header.Set("Authorization", "InvalidFormat")

	// Mock request with token missing Bearer prefix
	reqTokenNoBearer := httptest.NewRequest("GET", "/test", nil)
	reqTokenNoBearer.Header.Set("Authorization", "valid_token")

	// Mock request with invalid JSON in Secrets Manager response
	reqInvalidSecretJSON := httptest.NewRequest("GET", "/test", nil)
	reqInvalidSecretJSON.Header.Set("Authorization", "Bearer valid_token")

	// Mock request with JWT token validation failure
	reqInvalidToken := httptest.NewRequest("GET", "/test", nil)
	reqInvalidToken.Header.Set("Authorization", "Bearer invalid_token")

	// Mock request with valid JWT token and secret retrieval error
	reqValidTokenSecretError := httptest.NewRequest("GET", "/test", nil)
	reqValidTokenSecretError.Header.Set("Authorization", "Bearer valid_token")

	// Set environment variables for testing
	os.Setenv("REGION", "us-east-1")
	os.Setenv("SECRET", "your-secret-name")
	
	t.Run("Missing Authorization Header", func(t *testing.T) {
		rr := httptest.NewRecorder()
		ctx := context.Background()

		handler := Authenticate(mockHandler)
		handler.ServeHTTP(rr, reqMissingAuthHeader.WithContext(ctx))

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Status code does not match")
	})

	t.Run("Invalid Authorization Header Format", func(t *testing.T) {
		rr := httptest.NewRecorder()
		ctx := context.Background()

		handler := Authenticate(mockHandler)
		handler.ServeHTTP(rr, reqInvalidAuthHeader.WithContext(ctx))

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Status code does not match")
	})

	t.Run("Token Missing Bearer Prefix", func(t *testing.T) {
		rr := httptest.NewRecorder()
		ctx := context.Background()

		handler := Authenticate(mockHandler)
		handler.ServeHTTP(rr, reqTokenNoBearer.WithContext(ctx))

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Status code does not match")
	})

	t.Run("Valid Token with Missing Secret Environment Variables", func(t *testing.T) {
		// Unset environment variables for this test case
		os.Unsetenv("REGION")
		os.Unsetenv("SECRET")

		rr := httptest.NewRecorder()
		ctx := context.Background()

		handler := Authenticate(mockHandler)
		handler.ServeHTTP(rr, reqValidToken.WithContext(ctx))

		assert.Equal(t, http.StatusInternalServerError, rr.Code, "Status code does not match")

		// Reset environment variables
		os.Setenv("REGION", "us-east-1")
		os.Setenv("SECRET", "your-secret-name")
	})

	t.Run("Invalid Secret JSON", func(t *testing.T) {
		rr := httptest.NewRecorder()
		ctx := context.Background()

		handler := Authenticate(mockHandler)
		handler.ServeHTTP(rr, reqInvalidSecretJSON.WithContext(ctx))

		assert.Equal(t, http.StatusNotFound, rr.Code, "Status code does not match")
	})

	t.Run("Valid Token with Secret Retrieval Error", func(t *testing.T) {
		// Unset environment variables for this test case
		os.Unsetenv("REGION")
		os.Unsetenv("SECRET")

		rr := httptest.NewRecorder()
		ctx := context.Background()

		handler := Authenticate(mockHandler)
		handler.ServeHTTP(rr, reqValidTokenSecretError.WithContext(ctx))

		assert.Equal(t, http.StatusInternalServerError, rr.Code, "Status code does not match")
		// Reset environment variables
		os.Setenv("REGION", "us-east-1")
		os.Setenv("SECRET", "your-secret-name")
	})
}





// MockSecretsManagerClient is a mock of the SecretsManagerClient interface
type MockSecretsManagerClient struct {
	GetSecretValueFunc func(ctx context.Context, input *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

// GetSecretValue implements SecretsManagerClient
func (m *MockSecretsManagerClient) GetSecretValue(ctx context.Context, input *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	return m.GetSecretValueFunc(ctx, input, optFns...)
}

func TestGetSecret(t *testing.T) {
	originalSecretsManagerClient := secretsManagerClient
	defer func() { secretsManagerClient = originalSecretsManagerClient }()

	secretName := "testSecret"
	region := "us-west-2"

	validSecret := &models.SecretsManagerSecret{
		UserPoolID: "testPoolID",
		Region:     "us-west-2",
	}
	validSecretString, _ := json.Marshal(validSecret)

	secretsManagerClient = &MockSecretsManagerClient{
		GetSecretValueFunc: func(ctx context.Context, input *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
			if *input.SecretId == secretName {
				return &secretsmanager.GetSecretValueOutput{
					SecretString: aws.String(string(validSecretString)),
				}, nil
			}
			return nil, errors.New("some error")
		},
	}

	// Positive case: valid secret
	secret, err := GetSecretValue(region, secretName)
	assert.NoError(t, err)
	assert.NotNil(t, secret)
	assert.Equal(t, "testPoolID", secret.UserPoolID)
	assert.Equal(t, "us-west-2", secret.Region)

	// Negative case: secret retrieval error
	secretsManagerClient.(*MockSecretsManagerClient).GetSecretValueFunc = func(ctx context.Context, input *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
		return nil, errors.New("some error")
	}

	secret, err = GetSecretValue(region, secretName)
	assert.Error(t, err)
	assert.Nil(t, secret)
	assert.IsType(t, SecretRetrievalError{}, err)
	assert.Equal(t, "Secret retrieval error: some error", err.Error())

	// Negative case: secret string is nil
	secretsManagerClient.(*MockSecretsManagerClient).GetSecretValueFunc = func(ctx context.Context, input *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
		return &secretsmanager.GetSecretValueOutput{
			SecretString: nil,
		}, nil
	}

	secret, err = GetSecretValue(region, secretName)
	assert.Error(t, err)
	assert.Nil(t, secret)
	assert.IsType(t, SecretRetrievalError{}, err)
	assert.Equal(t, "Secret retrieval error: secret string is nil", err.Error())

	// Negative case: invalid JSON in secret string
	secretsManagerClient.(*MockSecretsManagerClient).GetSecretValueFunc = func(ctx context.Context, input *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
		return &secretsmanager.GetSecretValueOutput{
			SecretString: aws.String("invalid json"),
		}, nil
	}

	secret, err = GetSecretValue(region, secretName)
	assert.Error(t, err)
	assert.Nil(t, secret)
	assert.Contains(t, err.Error(), "invalid character")
}
