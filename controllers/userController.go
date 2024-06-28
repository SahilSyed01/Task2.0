package controllers

import (
	"context"
	"fmt"
	"go-chat-app/database"

	"github.com/go-playground/validator/v10"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type JwtAuthenticator func(ctx context.Context, region, userPoolID, tokenString string) (interface{}, error)

// MockValidateToken is a mock implementation of JwtAuthenticator for testing
//var MockValidateToken JwtAuthenticator

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	// if err != nil {
	//     log.Panic(err)
	// }
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}
	return check, msg
}
