package controllers

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"

    // "go-chat-app/database"
    // "go-chat-app/models"

    "github.com/ShreerajShettyK/cognitoJwtAuthenticator"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


func GetUsers(w http.ResponseWriter, r *http.Request) {
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
 
    // Token is valid, proceed with fetching users
    recordPerPage, err := strconv.Atoi(r.URL.Query().Get("recordPerPage"))
    if err != nil || recordPerPage < 1 {
        recordPerPage = 10 // Default value for recordPerPage
    }
 
    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil || page < 1 {
        page = 1 // Default value for page
    }
 
    startIndex := (page - 1) * recordPerPage
 
    matchStage := bson.D{{"$match", bson.D{{}}}}
    groupStage := bson.D{{"$group", bson.D{
        {"_id", bson.D{{"_id", "null"}}},
        {"total_count", bson.D{{"$sum", 1}}},
        {"data", bson.D{{"$push", bson.D{
            {"email", "$email"},
            {"first_name", "$first_name"},
            {"last_name", "$last_name"},
            {"password", "$password"},
            {"phone", "$phone"},
            {"user_id", "$user_id"},
        }}}},
    }}}
    projectStage := bson.D{
        {"$project", bson.D{
            {"_id", 0}, // Exclude the _id field
            {"total_count", 1},
            {"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
        }},
    }
 
    result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
        matchStage, groupStage, projectStage,
    })
    if err != nil {
        http.Error(w, "error occurred while listing user items", http.StatusInternalServerError)
        return
    }
 
    // Check if the response is empty
    if !result.Next(ctx) {
        http.Error(w, "No users found", http.StatusNotFound)
        return
    }
 
    // Custom struct for the response
    type UserResponse struct {
        TotalCount int      `json:"total_count"`
        UserItems  []bson.M `json:"user_items"`
    }
 
    // Decode the response into a temporary variable
    var tempResponse struct {
        TotalCount int      `bson:"total_count"`
        UserItems  []bson.M `bson:"user_items"`
    }
    if err := result.Decode(&tempResponse); err != nil {
        http.Error(w, fmt.Sprintf("error occurred while decoding user items: %v", err), http.StatusInternalServerError)
        return
    }
 
    // Convert the temporary response into the final UserResponse struct
    response := UserResponse{
        TotalCount: tempResponse.TotalCount,
        UserItems:  tempResponse.UserItems,
    }
 
    // Encode the custom response and send it
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}