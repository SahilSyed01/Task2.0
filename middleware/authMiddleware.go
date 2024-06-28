// middleware.go
package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ShreerajShettyK/cognitoJwtAuthenticator"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var (
	cognitovalidate = cognitoJwtAuthenticator.ValidateToken
	secretsClient   SecretsManagerClient
)

type SecretsManagerClient interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

type SecretsManagerSecret struct {
	UserPoolID string `json:"USER_POOL_ID"`
	Region     string `json:"REGION"`
}

func init() {
	// Initialize the secrets manager client
	region := os.Getenv("REGION")
	cfg,_ := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	// if err != nil {
	// 	log.Fatalf("Failed to load AWS config: %v", err)
	// }
	secretsClient = secretsmanager.NewFromConfig(cfg)
}

func Authenticate(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		authToken := strings.Split(authHeader, "Bearer ")
		if len(authToken) != 2 {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := authToken[1]

		// Fetch secrets from environment variables
		//region := os.Getenv("REGION")
		secretName := os.Getenv("SECRET")

		// Retrieve secret value using the initialized secretsClient
		input := &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secretName),
		}
		result, err := secretsClient.GetSecretValue(r.Context(), input)
		if err != nil {
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
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed with the request
		next.ServeHTTP(w, r)
	}
}
