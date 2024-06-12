package controllers

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strings"
    "time"

    "go-chat-app/helpers"
    "go-chat-app/models"

    "github.com/ShreerajShettyK/cognitoJwtAuthenticator"
    "go.mongodb.org/mongo-driver/bson"
)

func Login(w http.ResponseWriter, r *http.Request) {
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
 
    // Token is valid, proceed with login logic
 
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
    successMsg := map[string]string{"Success": "True", "ui_client_token": uiClientToken}
    response, err := json.Marshal(successMsg)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
 
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(response)
}
