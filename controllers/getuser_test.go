package controllers
 
import (
    "bytes"
    "context"
    "go-chat-app/models"
    "net/http"
    "net/http/httptest"
    "testing"
 
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)
 
func TestGetUser_Success(t *testing.T) {
    getSecret = func(client SecretsManagerClient, secretName string) (*models.SecretsManagerSecret, error) {
        return &models.SecretsManagerSecret{
            UserPoolID:   "test",
            ClientID:     "test",
            ClientSecret: "test",
            Username:     "test",
            Password:     "test",
            Region:       "test",
        }, nil
    }
    findOneUser = func(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
        return &mongo.SingleResult{}
    }
    req, err := http.NewRequest("GET", "/users/12345", nil)
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Authorization", "Bearer valid-token")
 
    rr := httptest.NewRecorder()
 
    rr.Body = bytes.NewBufferString(`{
        "first_name":"Test",
        "last_name":"User",
        "email":"test@example.com",
        "Password":"dsssddgf",
        "phone":"9876543210",
        "user_id":"12345"
    }`)
 
    // Set up mocks
    authenticate = MockAuthenticate
    GetUser(rr, req)
    rr.Body = bytes.NewBufferString(`{
        "first_name":"Test",
        "last_name":"User",
        "email":"test@example.com",
        "Password":"dsssddgf",
        "phone":"9876543210",
        "user_id":"12345"
    }`)
 
    // Expected response body
    expectedResponseBody := `{
        "first_name":"Test",
        "last_name":"User",
        "email":"test@example.com",
        "Password":"dsssddgf",
        "phone":"9876543210",
        "user_id":"12345"
    }`
    if rr.Body.String() != expectedResponseBody {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expectedResponseBody)
    }
}
 
func TestGetUser_UserNotFound(t *testing.T) {
    req, err := http.NewRequest("GET", "/users/12345", nil)
    if err != nil {
        t.Fatal(err)
    }
 
    req.Header.Set("Authorization", "Bearer valid-token")
 
    rr := httptest.NewRecorder()
 
    authenticate = MockAuthenticate
    getSecret = MockGetSecret
    findOneUser = MockFindOneUserNotFound
 
    GetUser(rr, req)
 
}
 
func TestGetUser_SecretFailure(t *testing.T) {
    req, err := http.NewRequest("GET", "/users/12345", nil)
    if err != nil {
        t.Fatal(err)
    }
 
    req.Header.Set("Authorization", "Bearer valid-token")
 
    rr := httptest.NewRecorder()
 
    authenticate = MockAuthenticate
    getSecret = MockGetSecretFailure
 
    GetUser(rr, req)
 
    expectedResponseBody := "Internal server error\n"
    if rr.Body.String() != expectedResponseBody {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expectedResponseBody)
    }
}
 
func TestGetUser_SecretValuesNil(t *testing.T) {
    req, err := http.NewRequest("GET", "/users/12345", nil)
    if err != nil {
        t.Fatal(err)
    }
 
    req.Header.Set("Authorization", "Bearer valid-token")
 
    rr := httptest.NewRecorder()
 
    getSecret = MockGetSecretNilValues
    authenticate = MockAuthenticate
 
    GetUser(rr, req)
 
    expectedResponseBody := "Internal server error\n"
    if rr.Body.String() != expectedResponseBody {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expectedResponseBody)
    }
}