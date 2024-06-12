package controllers

import (
	"context"
	"encoding/json"
	// "fmt"
	"log"
	"net/http"
	"time"

	"go-chat-app/middleware"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go-chat-app/models"
	"go-chat-app/helpers"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// Extract the JWT token and validate it using the middleware
	middleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Fetch secrets from Secrets Manager
		region := "us-east-1" // Set your AWS region here
		secretName := "myApp/mongo-db-credentials"
		secretResult, err := getSecret(region, secretName)
		if err != nil {
			log.Printf("Error fetching secret: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		secret := secretResult
		if secret == nil || secret.UserPoolID == nil || secret.Region == nil {
			log.Println("Secret, UserPoolID, or Region is nil")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Your login logic goes here
		var user models.User
		var foundUser models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			http.Error(w, "email or password is incorrect", http.StatusUnauthorized)
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !passwordIsValid {
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}

		// Generate token with First_name and UID
		token, err := helpers.GenerateToken(*foundUser.First_name, foundUser.User_id)
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
