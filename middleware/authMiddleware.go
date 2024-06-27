package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-chat-app/models"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ShreerajShettyK/cognitoJwtAuthenticator"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/smithy-go"
)

type SecretsManagerSecret struct {
	UserPoolID string `json:"USER_POOL_ID"`
	Region     string `json:"REGION"`
}

var (
	cognitovalidate = cognitoJwtAuthenticator.ValidateToken
)

func Authenticate(next http.Handler) http.HandlerFunc {
	//log.Println("insise auth func")
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the JWT token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Split the header value to extract the token part
		authToken := strings.Split(authHeader, "Bearer ")
		log.Println("auth token", authToken)
		if len(authToken) != 2 {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := authToken[1]

		// Fetch secrets from environment variables
		region := os.Getenv("REGION")
		secretName := os.Getenv("SECRET")

		// Get AWS config
		cfg, err := config.LoadDefaultConfig(r.Context(), config.WithRegion(region))
		if err != nil {
			log.Printf("Failed to load AWS config: %v", err)
			http.Error(w, "Failed to load AWS config", http.StatusInternalServerError)
			return
		}

		// Create Secrets Manager service client
		svc := secretsmanager.NewFromConfig(cfg)

		// Retrieve secret value
		input := &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secretName),
		}
		result, err := svc.GetSecretValue(r.Context(), input)
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && apiErr.ErrorCode() == "ResourceNotFoundException" {
				http.Error(w, "Secret not found", http.StatusNotFound)
				return
			}
			log.Printf("Error fetching secret: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Unmarshal secret string
		secret := &SecretsManagerSecret{}
		if err := json.Unmarshal([]byte(*result.SecretString), secret); err != nil {
			log.Printf("Failed to unmarshal secret string: %v", err)
			http.Error(w, "Failed to unmarshal secret", http.StatusInternalServerError)
			return
		}
		// Validate the JWT token
		_, err = cognitovalidate(r.Context(), secret.Region, secret.UserPoolID, tokenString)
		if err != nil {
			log.Printf("Token validation error: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized) // Set HTTP status code to 401
			return
		}

		// Token is valid, proceed with the request
		next.ServeHTTP(w, r)
	}
}

// SecretRetrievalError represents an error that occurred during secret retrieval.
type SecretRetrievalError struct {
	Message string
}

func (e SecretRetrievalError) Error() string {
	return fmt.Sprintf("Secret retrieval error: %s", e.Message)
}

// SecretsManagerClient is an interface for Secrets Manager client methods
type SecretsManagerClient interface {
	GetSecretValue(ctx context.Context, input *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

var secretsManagerClient SecretsManagerClient

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}
	secretsManagerClient = secretsmanager.NewFromConfig(cfg)
}

func GetSecretValue(region, secretName string) (*models.SecretsManagerSecret, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	}

	result, err := secretsManagerClient.GetSecretValue(context.Background(), input)
	if err != nil {
		return nil, SecretRetrievalError{Message: err.Error()}
	}

	if result.SecretString == nil {
		return nil, SecretRetrievalError{Message: "secret string is nil"}
	}

	secret := &models.SecretsManagerSecret{}
	err = json.Unmarshal([]byte(*result.SecretString), secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}
