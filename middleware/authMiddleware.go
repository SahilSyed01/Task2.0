package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ShreerajShettyK/cognitoJwtAuthenticator"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)
type SecretsManagerSecret struct {
	UserPoolID string `json:"USER_POOL_ID"`
	Region     string `json:"REGION"`
	
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the JWT token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Split the header value to extract the token part
		authToken := strings.Split(authHeader, "Bearer ")
		if len(authToken) != 2 {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		uiClientToken := authToken[1]

		// Fetch secrets from environment variables
		region := os.Getenv("REGION")
		secretName := os.Getenv("SECRET")
		secretResult, err := getSecretvalue(region, secretName)
		if err != nil {
			log.Printf("Error fetching secret: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		secret := secretResult.Secret
		if secret == nil || secret.UserPoolID == "" || secret.Region == "" {
			log.Println("Secret, UserPoolID, or Region is nil")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Validate the JWT token
		ctx := context.Background()
		tokenString := uiClientToken

		_, err = cognitoJwtAuthenticator.ValidateToken(ctx, secret.Region, secret.UserPoolID, tokenString)
		if err != nil {
			log.Printf("Token validation error: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed with the request
		next.ServeHTTP(w, r)
	})
}

type SecretResult struct {
	Secret *SecretsManagerSecret
	Err    error
}

func getSecretvalue(region string, secretName string) (SecretResult, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return SecretResult{Err: err}, err
	}

	svc := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		return SecretResult{Err: err}, err
	}

	if result.SecretString == nil {
		return SecretResult{Err: fmt.Errorf("secret string is nil")}, nil
	}

	secret := &SecretsManagerSecret{}
	err = json.Unmarshal([]byte(*result.SecretString), secret)
	if err != nil {
		return SecretResult{Err: err}, err
	}

	return SecretResult{Secret: secret}, nil
}

