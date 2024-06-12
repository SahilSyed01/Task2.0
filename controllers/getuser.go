package controllers

import (
	"context"
	"encoding/json"
	// "fmt"
	"log"
	"net/http"
	// "strings"
     "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go-chat-app/middleware"
	"go-chat-app/models"
	"time"

)

func GetUser(w http.ResponseWriter, r *http.Request) {
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

		// Token is valid, proceed with fetching the user
		userID := r.URL.Path[len("/users/"):]

		var user models.User
		err = userCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Define a custom response struct without the _id field
		type UserResponse struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Password  string `json:"Password"`
			Email     string `json:"email"`
			Phone     string `json:"phone"`
			UserID    string `json:"user_id"`
		}

		// Create a response object
		response := UserResponse{
			FirstName: *user.First_name,
			LastName:  *user.Last_name,
			Password:  *user.Password,
			Email:     *user.Email,
			Phone:     *user.Phone,
			UserID:    user.User_id,
		}

		// Encode the response object into JSON and send it
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})).ServeHTTP(w, r)
}
