package controllers

import (
	"bytes"
	"encoding/json"
	"go-chat-app/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)
func TestSignupWrongToken(t *testing.T) {
	// Prepare a request body with valid user data
	requestBody := []byte(`{"email": "test@example.com", "password": "password123", "phone": "1234567890"}`)
	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Set an incorrect or expired token in the Authorization header
	req.Header.Set("Authorization", "Bearer your-wrong-token")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the HTTP request with the Signup handler
	http.HandlerFunc(Signup).ServeHTTP(rr, req)

	// Check if the status code is Unauthorized (401)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestSignupMissingAuthHeader(t *testing.T) {
    req, err := http.NewRequest("POST", "/signup", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(Signup)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusUnauthorized {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
    }
}

// More negative test cases for Signup function can be added similarly...

// Negative test case for GetUsers function: Missing Authorization Header
func TestGetUsersMissingAuthHeader(t *testing.T) {
    req, err := http.NewRequest("GET", "/users", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(GetUsers)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusUnauthorized {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
    }
}

// More negative test cases for GetUsers function can be added similarly...

// Negative test case for GetUser function: Missing Authorization Header
func TestGetUserMissingAuthHeader(t *testing.T) {
    req, err := http.NewRequest("GET", "/users/userID", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(GetUser)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusUnauthorized {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
    }
}

// Negative test case for GetUsers function: Wrong Token
func TestGetUsersWrongToken(t *testing.T) {
    req, err := http.NewRequest("GET", "/users", nil)
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Authorization", "Bearer wrong_token_here")

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(GetUsers)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusUnauthorized {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
    }
}

// Negative test case for GetUser function: Wrong Token
func TestGetUserWrongToken(t *testing.T) {
    req, err := http.NewRequest("GET", "/users/userID", nil)
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Authorization", "Bearer wrong_token_here")

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(GetUser)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusUnauthorized {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
    }
}

func TestHashPassword(t *testing.T) {
    password := "testPassword"
    hashedPassword := HashPassword(password)

    if hashedPassword == "" {
        t.Error("Hashed password is empty")
    }
}

// Negative test case for VerifyPassword function: Incorrect Password
func TestVerifyPasswordIncorrectPassword(t *testing.T) {
    userPassword := "correctPassword"
    providedPassword := "incorrectPassword"

    check, msg := VerifyPassword(userPassword, providedPassword)

    if check {
        t.Error("Expected password verification to fail, but it passed")
    }
    if msg != "email or password is incorrect" {
        t.Errorf("Expected error message 'email or password is incorrect', but got '%s'", msg)
    }
}

// Positive test case for VerifyPassword function: Correct Password
func TestVerifyPasswordCorrectPassword(t *testing.T) {
    userPassword := "correctPassword"
    providedPassword := HashPassword("correctPassword")

    check, msg := VerifyPassword(userPassword, providedPassword)

    if !check {
        t.Error("Expected password verification to pass, but it failed")
    }
    if msg != "" {
        t.Errorf("Expected no error message, but got '%s'", msg)
    }
}
func TestLogin(t *testing.T) {
    t.Run("InvalidUser", func(t *testing.T) {
        // Negative test case: Invalid user credentials provided.
        // Create an invalid user object with incorrect email or password
        email := "invalid@example.com"
        password := "incorrectPassword"
        invalidUser := models.User{
            Email:    &email,
            Password: &password,
            // Add other required fields here...
        }

        // Marshal invalid user object to JSON
        requestBody, err := json.Marshal(invalidUser)
        assert.NoError(t, err)

        // Create a request with valid JSON body
        req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
        assert.NoError(t, err)

        // Create a response recorder
        rr := httptest.NewRecorder()

        // Call the Login handler function
        Login(rr, req)

        // Assert the response status code is 401 Unauthorized
        assert.Equal(t, http.StatusUnauthorized, rr.Code)
    })

    // Add more test cases for negative scenarios...
}
