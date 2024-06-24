package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	clients "go-chat-app/clients"
	"go-chat-app/helpers"
	"go-chat-app/middleware"
	"go-chat-app/models"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	getAWSConfig   = clients.GetAWSConfig
	getSMClient    = clients.GetSecretsManagerClient
	verifyPassword = VerifyPassword
	generateToken  = helpers.GenerateToken
	authenticate   = middleware.Authenticate
)

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var foundUser models.User
	// Ensure request body is not nil
	if r.Body == nil {
		http.Error(w, "Missing request body", http.StatusBadRequest)
		return
	}

	// Extract the JWT token and validate it using the middleware
	authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Fetch secrets from environment variables
		// region := os.Getenv("REGION")
		secretName := os.Getenv("SECRET")

		awsConfig, err := getAWSConfig()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		smClient := getSMClient(awsConfig)
		secretResult, err := getSecret(smClient, secretName)
		if err != nil {
			log.Printf("Error fetching secret: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		secret := secretResult
		if secret == nil || secret.UserPoolID == "" || secret.Region == "" {
			log.Println("Secret, UserPoolID, or Region is nil")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Your existing login logic goes here...

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = findOneUser(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			http.Error(w, "email or password is incorrect", http.StatusUnauthorized)
			return
		}

		passwordIsValid, msg := verifyPassword(user.Password, foundUser.Password)
		if !passwordIsValid {
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}

		// Generate token with First_name and UID
		token, err := generateToken(foundUser.First_name, foundUser.User_id)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// Set token in response header
		w.Header().Set("Authorization", "Bearer "+token)

		// Respond with a simple success message in JSON format
		successMsg := map[string]string{"Success": "True"}
		response, err := json.Marshal(successMsg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	})).ServeHTTP(w, r)
}
