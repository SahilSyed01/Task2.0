package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"fmt"
     "go-chat-app/middleware"
	// "github.com/ShreerajShettyK/cognitoJwtAuthenticator"
	"go.mongodb.org/mongo-driver/bson"
	"go-chat-app/models"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	// Extract the JWT token and validate it using the middleware
	middleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Fetch secrets from Secrets Manager
		region := "us-east-1" // Set your AWS region here
		secretName := "myApp/mongo-db-credentials"
		secretResult,err := getSecret(region, secretName)
		if err != nil {
			log.Printf("Error fetching secret: %v",err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		secret := secretResult
		if secret == nil || secret.UserPoolID == nil || secret.Region == nil {
			log.Println("Secret, UserPoolID, or Region is nil")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Your signup logic goes here
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			http.Error(w, validationErr.Error(), http.StatusBadRequest)
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			http.Error(w, "error occurred while checking for the email", http.StatusInternalServerError)
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			http.Error(w, "error occurred while checking for the phone number", http.StatusInternalServerError)
			return
		}

		if count > 0 {
			http.Error(w, "this email or phone number already exists", http.StatusInternalServerError)
			return
		}

		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(resultInsertionNumber)
	})).ServeHTTP(w, r)
}
