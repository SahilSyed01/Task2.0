package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/joho/godotenv"

	clients "go-chat-app/clients"
	"go-chat-app/cognito"
	"go-chat-app/controllers"
	"go-chat-app/routes"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Load AWS configuration
	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	cogitoClient := clients.GetCognitoClient(awsConfig)
	smClient := clients.GetSecretsManagerClient(awsConfig)

	// Setup routes
	routes.AuthRoutes()
	routes.UserRoutes()

	// New route to get JWT token
	http.HandleFunc("/get-jwt-token", func(w http.ResponseWriter, r *http.Request) {
		// Fetch the secret from Secrets Manager
		secretName := os.Getenv("SECRET")
		secret, err := controllers.GetSecret(smClient, secretName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Retrieve values from the secret
		userPoolID := secret.UserPoolID
		clientID := secret.ClientID
		clientSecret := secret.ClientSecret
		username := secret.Username
		password := secret.Password

		token, err := cognito.GetJWTToken(cogitoClient, userPoolID, clientID, clientSecret, username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"jwt_token": token})
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
