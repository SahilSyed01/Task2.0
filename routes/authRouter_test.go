package routes

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestAuthRoutes(t *testing.T) {
    // Create a request for the "/users/signup" route
    reqSignup, err := http.NewRequest("GET", "/users/signup", nil)
    if err != nil {
        t.Fatal(err)
    }

    // Create a request for the "/users/login" route
    reqLogin, err := http.NewRequest("GET", "/users/login", nil)
    if err != nil {
        t.Fatal(err)
    }

    // Create a ResponseRecorder to record the response
    rr := httptest.NewRecorder()

    // Call the AuthRoutes function to set up the routes
    AuthRoutes()

    // Serve the "/users/signup" route
    http.DefaultServeMux.ServeHTTP(rr, reqSignup)

    // Check if the status code is what you expect (401 Unauthorized)
    if status := rr.Code; status != http.StatusUnauthorized {
        t.Errorf("handler returned wrong status code for /users/signup: got %v want %v",
            status, http.StatusUnauthorized)
    }

    // Serve the "/users/login" route
    http.DefaultServeMux.ServeHTTP(rr, reqLogin)

    // Check if the status code is what you expect (401 Unauthorized)
    if status := rr.Code; status != http.StatusUnauthorized {
        t.Errorf("handler returned wrong status code for /users/login: got %v want %v",
            status, http.StatusUnauthorized)
    }
}
