package main

import (
	"log"
	"net/http"
	"user-management-service/config"
	"user-management-service/handlers"
	"user-management-service/repository"
	"user-management-service/utils"
	"github.com/gorilla/mux"
)

func main() {
	config.LoadConfig()
	err := repository.Connect()
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v", err)
	}
	utils.InitCognito()

	r := mux.NewRouter()
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	http.Handle("/", r)
	log.Println("Starting server on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
