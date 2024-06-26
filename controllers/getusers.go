package controllers
 
import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"
 
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)
 
var (
    aggregate = userCollection.Aggregate
)
func GetUsers(w http.ResponseWriter, r *http.Request) {
    // Extract the JWT token and validate it using the middleware
    authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()
 
        // Fetch secrets from environment variables
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
        cursor, err := aggregate(ctx, mongo.Pipeline{
            matchStage, groupStage, projectStage,
        })
        if err != nil {
            log.Printf("Error in MongoDB aggregation: %v", err)
            http.Error(w, "Error occurred while listing user items", http.StatusInternalServerError)
            return
        }
        defer cursor.Close(ctx)
 
        // Check if no results found
        if !cursor.Next(ctx) {
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
        if err := cursor.Decode(&tempResponse); err != nil {
            log.Printf("Error decoding MongoDB cursor: %v", err)
            http.Error(w, "Error occurred while decoding user items", http.StatusInternalServerError)
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
            log.Printf("Error encoding JSON response: %v", err)
            http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
            return
        }
    })).ServeHTTP(w, r)
}