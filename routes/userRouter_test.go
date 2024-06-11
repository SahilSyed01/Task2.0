package routes

import (
    // "go-chat-app/controllers"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestUserRoutes(t *testing.T) {
    // Setup the routes
    UserRoutes()

    tests := []struct {
        method   string
        url      string
        expected int
    }{
        {"GET", "/users", http.StatusUnauthorized},
        {"GET", "/users/1", http.StatusUnauthorized},
    }

    for _, tt := range tests {
        req, err := http.NewRequest(tt.method, tt.url, nil)
        if err != nil {
            t.Fatal(err)
        }

        rr := httptest.NewRecorder()
        http.DefaultServeMux.ServeHTTP(rr, req)

        if status := rr.Code; status != tt.expected {
            t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expected)
        }
    }
}
