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
    userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
    jwtParseWithClaim                = jwt.ParseWithClaims
    secretKey        string
    newWithClaimsFunc = jwt.NewWithClaims
    signTokenFunc     = func(token *jwt.Token, key interface{}) (string, error) {
        return token.SignedString(key)
    }
)
 
func SetSecretKey(key string) {
    secretKey = key
}
 
func SetNewWithClaimsFunc(f func(jwt.SigningMethod, jwt.Claims) *jwt.Token) {
    newWithClaimsFunc = f
}
 
func SetSignTokenFunc(f func(*jwt.Token, interface{}) (string, error)) {
    signTokenFunc = f
}
 
func GenerateToken(firstName string, userID string) (string, error) {
    var signedToken string
    accessTokenClaims := &SignedDetails{
        First_name: firstName,
        Uid:        userID,
    }
 
    accessToken := newWithClaimsFunc(jwt.SigningMethodHS256, accessTokenClaims)
    signedToken, err := signTokenFunc(accessToken, []byte(secretKey))
    if err != nil {
        log.Println("Error generating access token:", err)
        return signedToken, err
    }
 
    return signedToken, err
}
 
func ValidateToken(signedToken string) (*SignedDetails, error) {
    var claims *SignedDetails
    token, err := jwtParseWithClaim(
        signedToken,
        &SignedDetails{},
        func(token *jwt.Token) (interface{}, error) {
            var err error
            return []byte(secretKey), err
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