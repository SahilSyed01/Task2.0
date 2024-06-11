package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	//"go-chat-app/controllers"
	"go-chat-app/routes"
)

func TestAuthRoutes(t *testing.T) {
	// Create a new HTTP request multiplexer
	mux := http.NewServeMux()

	// Register routes using AuthRoutes
	routes.AuthRoutes()

	// Create a test server with the created multiplexer
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Test signup route
	resp, err := http.Post(ts.URL+"/users/signup", "application/json", nil)
	if err != nil {
		t.Errorf("Error sending POST request to /users/signup: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Test login route
	resp, err = http.Post(ts.URL+"/users/login", "application/json", nil)
	if err != nil {
		t.Errorf("Error sending POST request to /users/login: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
