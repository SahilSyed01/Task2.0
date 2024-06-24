package helpers

import (
	"go-chat-app/database"
	"log"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
)



type SignedDetails struct {
	First_name string `json:"first_name"`
	Uid        string `json:"uid"`
	jwt.StandardClaims
}

var (
	userCollection    *mongo.Collection = database.OpenCollection(database.Client, "user")
	jwtParseWithClaim                   = jwt.ParseWithClaims
	
)

var SECRET_KEY string

func GenerateToken(firstName string, userID string) (string, error) {
	var signedToken string
	accessTokenClaims := &SignedDetails{
		First_name: firstName,
		Uid:        userID,
	}

	// Generate the access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	signedToken, err := accessToken.SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Println("Error generating access token:", err)
		return signedToken, err // Here, return the error
	}

	return signedToken, err // Return the signed token
}

func ValidateToken(signedToken string) (*SignedDetails, error) {
	var claims *SignedDetails
	token, err := jwtParseWithClaim(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			var err error
			return []byte(SECRET_KEY), err
		},
	)

	if err != nil {
		log.Println("Error parsing token:", err)
		return claims, err
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		log.Println("Error casting token claims")
		return claims, err
	}

	return claims, err
}
