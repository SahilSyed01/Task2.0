package controllers

import (
	"bytes"

	// "encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	//"time"

	"go-chat-app/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	//"github.com/stretchr/testify/assert"
)

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
		authenticate = func(next http.Handler) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				r.Header = http.Header{
					"Content-Type":  {"application/json"},
					"Authorization": {"Bearer asd.asd.asd"},
				}
			}
		}

		generateToken = func(firstName, userID string) (string, error) {
			return "test", nil
		}
		verifyPassword = func(userPassword, providedPassword string) (bool, string) {
			return true, ""
		}
		getAWSConfig = func() (aws.Config, error) {
			return aws.Config{
				Region: "test",
			}, nil
		}
		getSMClient = func(config aws.Config) *secretsmanager.Client {
			return &secretsmanager.Client{}
		}
		getSecret = func(client SecretsManagerClient, secretName string) (*models.SecretsManagerSecret, error) {
			return &models.SecretsManagerSecret{
				UserPoolID: "test",
				Region:     "test",
			}, nil
		}

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

func TestLogin_UnauthorizedIncorrectCredentials(t *testing.T) {
	// Prepare a request body with incorrect credentials
	requestBody := []byte(`{"email": "test@example.com", "password": "wrongpassword"}`)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock dependencies
	authenticate = func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {}
	}
	verifyPassword = func(userPassword, providedPassword string) (bool, string) {
		return false, ""
	}
	getAWSConfig = func() (aws.Config, error) {
		return aws.Config{
			Region: "test",
		}, nil
	}
	getSMClient = func(config aws.Config) *secretsmanager.Client {
		return &secretsmanager.Client{}
	}
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

	// Call the Login handler function directly
	Login(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	// Check the response body
	expectedResponseBody := "email or password is incorrect\n"
	if rr.Body.String() != expectedResponseBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedResponseBody)
	}
}

// Mocks for dependencies
