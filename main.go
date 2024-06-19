package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
    "github.com/joho/godotenv"

    "go-chat-app/cognito"
    "go-chat-app/controllers" 
    "go-chat-app/routes"
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

    // Load AWS configuration
    cfg, err := config.LoadDefaultConfig(context.Background())
    if err != nil {
        log.Fatalf("failed to load AWS config: %v", err)
    }

    // Create Cognito Identity Provider client
    svc := cognitoidentityprovider.NewFromConfig(cfg)

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

    // New route to get JWT token
    http.HandleFunc("/get-jwt-token", func(w http.ResponseWriter, r *http.Request) {
        // Fetch the secret from Secrets Manager
        secretName := os.Getenv("SECRET")
        region := os.Getenv("REGION")

        secret, err := controllers.GetSecret(region, secretName)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Retrieve values from the secret
        userPoolID := secret.UserPoolID
        clientID := secret.ClientID
        clientSecret := secret.ClientSecret
        username := secret.Username
        password := secret.Password

        token, err := cognito.GetJWTToken(svc, userPoolID, clientID, clientSecret, username, password)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(map[string]string{"jwt_token": token})
    })

    log.Fatal(http.ListenAndServe(":"+port, nil))
}
