package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ShreerajShettyK/cognitoJwtAuthenticator"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

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

		// Fetch secrets from Secrets Manager
		region := "us-east-1" // Set your AWS region here
		secretName := "myApp/mongo-db-credentials"
		secret, err := getSecret(region, secretName)
		if err != nil {
			log.Printf("Error fetching secret: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if secret == nil || secret.UserPoolID == nil || secret.Region == nil {
			log.Println("Secret, UserPoolID, or Region is nil")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Validate the JWT token
		ctx := context.Background()
		tokenString := uiClientToken

		_, err = cognitoJwtAuthenticator.ValidateToken(ctx, *secret.Region, *secret.UserPoolID, tokenString)
		if err != nil {
			log.Printf("Token validation error: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed with the request
		next.ServeHTTP(w, r)
	})
}

func getSecret(region, secretName string) (*SecretsManagerSecret, error) {
	sess := session.Must(session.NewSession())
	svc := secretsmanager.New(sess, &aws.Config{Region: aws.String(region)})

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return nil, err
	}

	if result.SecretString == nil {
		return nil, fmt.Errorf("secret string is nil")
	}

	secret := &SecretsManagerSecret{}
	err = json.Unmarshal([]byte(*result.SecretString), secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

type SecretsManagerSecret struct {
	UserPoolID *string `json:"USER_POOL_ID"`
	Region     *string `json:"REGION"`
	// Add other fields from your secret here
}
