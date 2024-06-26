package controllers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-chat-app/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mock functions
func MockAuthenticate(next http.Handler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        next.ServeHTTP(w, r)
    }
}
 
func MockGetAWSConfig() (aws.Config, error) {
    return aws.Config{
        Region: "test",
    }, nil
}
 
func MockGetSMClient(config aws.Config) *secretsmanager.Client {
    return &secretsmanager.Client{}
}
 
func MockGetSecret(client SecretsManagerClient, secretName string) (*models.SecretsManagerSecret, error) {
    return &models.SecretsManagerSecret{
        UserPoolID: "test",
        Region:     "test",
    }, nil
}
 
func MockVerifyPassword(userPassword, providedPassword string) (bool, string) {
    return true, ""
}
 
func MockGenerateToken(firstName, userID string) (string, error) {
    return "mocked_token", nil
}
 
func MockFindOneUser(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
    user := models.User{
        Email:      "test@example.com",
        Password:   "password123",
        First_name: "Test",
        User_id:    "12345",
    }
    doc, _ := bson.Marshal(user)
    return mongo.NewSingleResultFromDocument(doc, nil, nil)
}
func MockVerifyPasswordFailure(userPassword, providedPassword string) (bool, string) {
    return false, "password incorrect"
}

func TestLogin_Success(t *testing.T) {
    t.Run("test for successful login", func(t *testing.T) {
        requestBody := []byte(`{"email": "test@example.com", "password": "password123"}`)
        req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
        if err != nil {
            t.Fatal(err)
        }
        req.Header = http.Header{
            "Content-Type": {"application/json"},
        }
 
        // Create a ResponseRecorder to record the response
        rr := httptest.NewRecorder()
 
        // Mock dependencies
        authenticate = MockAuthenticate
        getAWSConfig = MockGetAWSConfig
        getSMClient = MockGetSMClient
        getSecret = MockGetSecret
        verifyPassword = MockVerifyPassword
        generateToken = MockGenerateToken
        findOneUser = MockFindOneUser
 
        // Call the Login handler function directly
        Login(rr, req)
 
        // Check the status code
        if status := rr.Code; status != http.StatusOK {
            t.Errorf("handler returned wrong status code: got %v want %v",
                status, http.StatusOK)
        }
 
        // Check the response body
        expectedResponseBody := `{"Success":"True"}`
        if rr.Body.String() != expectedResponseBody {
            t.Errorf("handler returned unexpected body: got %v want %v",
                rr.Body.String(), expectedResponseBody)
        }
 
        // Check the Authorization header
        expectedToken := "Bearer mocked_token"
        if rr.Header().Get("Authorization") != expectedToken {
            t.Errorf("handler returned unexpected Authorization header: got %v want %v",
                rr.Header().Get("Authorization"), expectedToken)
        }
    })
}
 
func TestLogin_BadRequestMissingBody(t *testing.T) {
    req, err := http.NewRequest("POST", "/login", nil)
    if err != nil {
        t.Fatal(err)
    }
 
    // Create a ResponseRecorder to record the response
    rr := httptest.NewRecorder()
 
    // Call the Login handler function directly
    Login(rr, req)
 
    // Check the status code
    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusBadRequest)
    }
}

func MockFindOneUserNotFound(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
    return mongo.NewSingleResultFromDocument(nil, nil, nil) // Simulate user not found
}

func TestLogin_UserNotFound(t *testing.T) {
    t.Run("test for user not found in MongoDB", func(t *testing.T) {
        requestBody := []byte(`{"email": "test@example.com", "password": "password123"}`)
        req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
        if err != nil {
            t.Fatal(err)
        }
        req.Header = http.Header{
            "Content-Type": {"application/json"},
        }

        // Create a ResponseRecorder to record the response
        rr := httptest.NewRecorder()

        // Mock dependency
        findOneUser = MockFindOneUserNotFound

        // Call the Login handler function directly
        Login(rr, req)

        // Check the status code
        if status := rr.Code; status != http.StatusUnauthorized {
            t.Errorf("handler returned wrong status code: got %v want %v",
                status, http.StatusUnauthorized)
        }
    })
}


func TestLogin_PasswordVerificationFailure(t *testing.T) {
    t.Run("test for password verification failure", func(t *testing.T) {
        requestBody := []byte(`{"email": "test@example.com", "password": "password123"}`)
        req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
        if err != nil {
            t.Fatal(err)
        }
        req.Header = http.Header{
            "Content-Type": {"application/json"},
        }

        // Create a ResponseRecorder to record the response
        rr := httptest.NewRecorder()

        // Mock dependency
        verifyPassword = MockVerifyPasswordFailure

        // Call the Login handler function directly
        Login(rr, req)

        // Check the status code
        if status := rr.Code; status != http.StatusUnauthorized {
            t.Errorf("handler returned wrong status code: got %v want %v",
                status, http.StatusUnauthorized)
        }
    })
}

func MockGetAWSConfigFailure() (aws.Config, error) {
    return aws.Config{}, errors.New("failed to fetch AWS config")
}

func TestLogin_FailureFetchAWSConfig(t *testing.T) {
    req, err := http.NewRequest("POST", "/login", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()

    // Mock dependency to simulate failure to fetch AWS config
    getAWSConfig = MockGetAWSConfigFailure

    Login(rr, req)

    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusBadRequest)
    }
    // Check the response body for the error message
    expectedResponseBody := "Missing request body"
    actualResponseBody := rr.Body.String()
    if actualResponseBody != expectedResponseBody+"\n" { // Add newline character here
        t.Errorf("handler returned unexpected body: got %v want %v",
            actualResponseBody, expectedResponseBody)
    }
}

// Mock the Authenticate middleware to simulate a successful authentication

// Mock function to simulate a failure in fetching secrets
func MockGetSecretFailure(client SecretsManagerClient, secretName string) (*models.SecretsManagerSecret, error) {
    return nil, errors.New("simulated secret retrieval failure")
}

func TestLogin_FailureFetchSecrets(t *testing.T) {
    requestBody := []byte(`{"email": "test@example.com", "password": "password123"}`)
    req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
    if err != nil {
        t.Fatal(err)
    }
    req.Header = http.Header{
        "Content-Type":   {"application/json"},
        "Authorization": {"Bearer valid-token"}, // Add a valid token for the test
    }

    rr := httptest.NewRecorder()

    // Mock dependencies
    authenticate = MockAuthenticate
    getAWSConfig = MockGetAWSConfig
    getSMClient = MockGetSMClient
    getSecret = MockGetSecretFailure

    Login(rr, req)

    if status := rr.Code; status != http.StatusInternalServerError {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusInternalServerError)
    }
    // Check the response body for the error message
    expectedResponseBody := "Internal server error\n" // Including newline character
    if rr.Body.String() != expectedResponseBody {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expectedResponseBody)
    }
}

func TestLogin_FailureDecodeJSON(t *testing.T) {
    t.Run("test for failure to decode JSON request body", func(t *testing.T) {
        requestBody := []byte(`{"email": "test@example.com", "password": "password123"`)
        req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody)) // Malformed JSON
        if err != nil {
            t.Fatal(err)
        }
        req.Header = http.Header{
            "Content-Type": {"application/json"},
        }

        // Create a ResponseRecorder to record the response
        rr := httptest.NewRecorder()

        // Mock dependencies
        authenticate = MockAuthenticate
        getAWSConfig = MockGetAWSConfig
        getSMClient = MockGetSMClient
        getSecret = MockGetSecret
        verifyPassword = MockVerifyPassword
        generateToken = MockGenerateToken
        findOneUser = MockFindOneUser

        // Call the Login handler function directly
        Login(rr, req)

        // Check the status code
        if status := rr.Code; status != http.StatusBadRequest {
            t.Errorf("handler returned wrong status code: got %v want %v",
                status, http.StatusBadRequest)
        }
    })
}


func MockGenerateTokenFailure(firstName, userID string) (string, error) {
    return "", errors.New("simulated token generation failure")
}
func TestLogin_FailureGenerateToken(t *testing.T) {
    t.Run("test for failure to generate token", func(t *testing.T) {
        requestBody := []byte(`{"email": "test@example.com", "password": "password123"}`)
        req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
        if err != nil {
            t.Fatal(err)
        }
        req.Header = http.Header{
            "Content-Type": {"application/json"},
        }

        // Create a ResponseRecorder to record the response
        rr := httptest.NewRecorder()

        // Mock dependencies
        authenticate = MockAuthenticate
        getAWSConfig = MockGetAWSConfig
        getSMClient = MockGetSMClient
        getSecret = MockGetSecret
        verifyPassword = MockVerifyPassword
        generateToken = MockGenerateTokenFailure
        findOneUser = MockFindOneUser

        // Call the Login handler function directly
        Login(rr, req)

        // Check the status code
        if status := rr.Code; status != http.StatusInternalServerError {
            t.Errorf("handler returned wrong status code: got %v want %v",
                status, http.StatusInternalServerError)
        }
    })
}

func TestLogin_IncorrectPassword(t *testing.T) {
    t.Run("test for incorrect password", func(t *testing.T) {
        requestBody := []byte(`{"email": "test@example.com", "password": "wrongpassword"}`)
        req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
        if err != nil {
            t.Fatal(err)
        }
        req.Header = http.Header{
            "Content-Type": {"application/json"},
        }

        // Create a ResponseRecorder to record the response
        rr := httptest.NewRecorder()

        // Mock dependencies
        authenticate = MockAuthenticate
        getAWSConfig = MockGetAWSConfig
        getSMClient = MockGetSMClient
        getSecret = MockGetSecret
        verifyPassword = MockVerifyPasswordFailure
        generateToken = MockGenerateToken
        findOneUser = MockFindOneUser

        // Call the Login handler function directly
        Login(rr, req)

        // Check the status code
        if status := rr.Code; status != http.StatusUnauthorized {
            t.Errorf("handler returned wrong status code: got %v want %v",
                status, http.StatusUnauthorized)
        }
    })
}
func MockGetSecretNilValues(client SecretsManagerClient, secretName string) (*models.SecretsManagerSecret, error) {
    return &models.SecretsManagerSecret{
        UserPoolID: "",
        Region:     "",
    }, nil
}

func TestLogin_SecretValuesNil(t *testing.T) {
    t.Run("test for secret values being nil", func(t *testing.T) {
        requestBody := []byte(`{"email": "test@example.com", "password": "password123"}`)
        req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
        if err != nil {
            t.Fatal(err)
        }
        req.Header = http.Header{
            "Content-Type": {"application/json"},
            "Authorization": {"Bearer testtoken"},
        }

        // Create a ResponseRecorder to record the response
        rr := httptest.NewRecorder()

        // Mock dependencies
        getSecret = MockGetSecretNilValues
        getAWSConfig = MockGetAWSConfig
        authenticate = MockAuthenticate

        // Call the Login handler function directly
        Login(rr, req)

        // Check the status code
        if status := rr.Code; status != http.StatusInternalServerError {
            t.Errorf("handler returned wrong status code: got %v want %v",
                status, http.StatusInternalServerError)
        }

        // Check the response body
        expectedResponseBody := "Internal server error\n"
        if rr.Body.String() != expectedResponseBody {
            t.Errorf("handler returned unexpected body: got %v want %v",
                rr.Body.String(), expectedResponseBody)
        }
    })
}
