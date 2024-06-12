package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "strconv"
	"strings"
	"time"

	// "go-chat-app/database"
	// "go-chat-app/helpers"
	"go-chat-app/models"

	"github.com/ShreerajShettyK/cognitoJwtAuthenticator"
	// "github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, r *http.Request) {
    var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
    defer cancel()
 
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
 
    // Validate the JWT token
    ctx = context.Background()
    tokenString := uiClientToken
 
    _, err = cognitoJwtAuthenticator.ValidateToken(ctx, *secret.Region, *secret.UserPoolID, tokenString)
    if err != nil {
        http.Error(w, fmt.Sprintf("Token validation error: %s", err), http.StatusUnauthorized)
        return
    }
 
    // Token is valid, proceed with signup logic
 
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
}
