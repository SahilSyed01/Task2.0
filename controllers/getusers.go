package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go-chat-app/middleware"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	// Extract the JWT token and validate it using the middleware
	middleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Fetch secrets from environment variables
		region := os.Getenv("REGION")
		secretName := os.Getenv("SECRET")
		secretResult, err := GetSecret(region, secretName)
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

		// Your get users logic goes here
		// Parse URL query parameters for pagination
		recordPerPage, err := strconv.Atoi(r.URL.Query().Get("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10 // Default value for recordPerPage
		}

		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 1 {
			page = 1 // Default value for page
		}

		startIndex := (page - 1) * recordPerPage

		// MongoDB aggregation pipeline stages
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

		// Aggregate pipeline execution
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
	})).ServeHTTP(w, r)
}
