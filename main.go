package main

import (
	"log"
	"net/http"
	"os"

	"go-chat-app/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Setup routes
	routes.AuthRoutes()
	routes.UserRoutes()

	http.HandleFunc("/api-1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":"Access granted for api-1"}`))
	})

	http.HandleFunc("/api-2", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":"Access granted for api-2"}`))
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
