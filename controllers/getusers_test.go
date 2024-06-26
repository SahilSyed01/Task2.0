package controllers
 
import (
    "context"
    "errors"
    "go-chat-app/models"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
 
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)
 
func MockAuthenticate2(next http.Handler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Mock token validation or authorization logic
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
            return
        }
 
        // Example: Validate the token format and value
        expectedToken := "Bearer valid-token"
        if authHeader != expectedToken {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }
 
        // Proceed with the request if token is valid
        next.ServeHTTP(w, r)
    }
}
 
func MockGetSecret2(client SecretsManagerClient, secretName string) (*models.SecretsManagerSecret, error) {
    return &models.SecretsManagerSecret{
        UserPoolID:   "test",
        ClientID:     "test",
        ClientSecret: "test",
        Username:     "test",
        Password:     "test",
        Region:       "test",
    }, nil
}
 
func TestGetUsers_Success(t *testing.T) {
    // Mock environment variables
    os.Setenv("SECRET", "test_secret") // Replace with your actual secret name
    defer os.Unsetenv("SECRET")
 
    // Mock authentication middleware function
    authenticate = MockAuthenticate2
 
    // Mock secrets manager client function
    getSecret = MockGetSecret2
 
    // Mock MongoDB aggregation function
    aggregate = func(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
        return nil, errors.New("MongoDB aggregation error")
    }
 
    // Create a mock request
    req, err := http.NewRequest("GET", "/users?recordPerPage=10&page=1", nil)
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Authorization", "Bearer valid-token") // Set valid token in header
 
    // Create a mock response recorder
    rr := httptest.NewRecorder()
 
    // Call the handler function
    GetUsers(rr, req)
 
    // Check the status code
    if status := rr.Code; status != http.StatusInternalServerError {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusInternalServerError)
    }
 
    // Check the content type
    expectedContentType := "text/plain; charset=utf-8"
    if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
        t.Errorf("handler returned wrong content type: got %v want %v",
            contentType, expectedContentType)
    }
}