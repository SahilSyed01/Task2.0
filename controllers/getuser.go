package controllers

import (
	"context"
	"encoding/json"
	"go-chat-app/models"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	getSecret   = GetSecret
	findOneUser = userCollection.FindOne
	// userCollection.FindOne
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Extract the JWT token and validate it using the middleware
	authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Fetch secrets from Secrets Manager
		// region := os.Getenv("REGION")
		secretName := os.Getenv("SECRET")
		secretResult, err := getSecret(secretsManagerClient, secretName)
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

		// Token is valid, proceed with fetching the user
		userID := r.URL.Path[len("/users/"):]

		var user models.User
		err = findOneUser(ctx, bson.M{"user_id": userID}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
	 json.NewEncoder(w).Encode(user);
	 
	})).ServeHTTP(w, r)
}
