package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"go-chat-app/models"

	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	countDocs     = userCollection.CountDocuments
	hashPassword  = HashPassword
	insertOneUser = userCollection.InsertOne
)

func Signup(w http.ResponseWriter, r *http.Request) {
	// Extract the JWT token and validate it using the middleware
	authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		awsConfig, err := getAWSConfig()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		smClient := getSMClient(awsConfig)
		// Fetch secrets from environment variables
		// region := os.Getenv("REGION")
		secretName := os.Getenv("SECRET")
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

		// Decode the request body into a user struct
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validate user input (e.g., using a validator library)
		validationErr := validate.Struct(user)
		if validationErr != nil {
			http.Error(w, validationErr.Error(), http.StatusBadRequest)
			return
		}

		// Check if the email already exists
		count, err := countDocs(r.Context(), bson.M{"phone": user.Phone})
		if err != nil {
			log.Println("Error counting documents:", err)
			http.Error(w, "error occurred while checking for the phone number", http.StatusInternalServerError)
			return
		}
		if count > 0 {
			http.Error(w, "this phone number already exists", http.StatusBadRequest)
			return
		}

		// Check if the phone number already exists
		count, err = countDocs(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			//log.Panic(err)
			http.Error(w, "error occurred while checking for the phone number", http.StatusInternalServerError)
			return
		}
		if count > 0 {
			http.Error(w, "this phone number already exists", http.StatusBadRequest)
			return
		}

		// Hash the user's password (you need to implement HashPassword function)
		password := hashPassword(user.Password)
		user.Password = password

		// Insert the user into MongoDB
		resultInsertionNumber, insertErr := insertOneUser(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		// Return a success response
		json.NewEncoder(w).Encode(resultInsertionNumber)
	})).ServeHTTP(w, r)
}

