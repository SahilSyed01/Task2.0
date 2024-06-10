package main

import (
	"fmt"
	// "strings"
	// "os"
	"go-chat-app/cognito"
)

// "log"
// "net/http"
// "os"

// "go-chat-app/routes"
// "github.com/joho/godotenv"

func main() {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8000"
	// }

	// // Setup routes
	// routes.AuthRoutes()
	// routes.UserRoutes()

	// http.HandleFunc("/api-1", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte(`{"success":"Access granted for api-1"}`))
	// })

	// http.HandleFunc("/api-2", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte(`{"success":"Access granted for api-2"}`))
	// })

	// log.Fatal(http.ListenAndServe(":"+port, nil))

	// cognitoDomain := "pujitha.auth.us-east-1.amazoncognito.com"

	// // Split the domain by "." and take the first part
	// parts := strings.Split(cognitoDomain, ".")

	client := &cognito.CognitoClient{
		ClientID:     "5hi3p0d0lvp7fcl1o05fchj8ui",
		ClientSecret: "1fe806n734kov79tm5h504fsf3h38v0nmd01e205kbuqaorfm5qn",
		// PoolID:       "us-east-1_bLcLm4KQ2",
		Region:       "us-east-1",
		Domain:       "mytestcognt", 
	}

	username := "mytestuser"
	password := "saipujitha"

	jwt, err := client.GetJWT(username, password)
	if err != nil {
		fmt.Println("Error getting JWT:", err)
		return
	}

	fmt.Println("JWT Token:", jwt)
}
