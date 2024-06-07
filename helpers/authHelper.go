package helpers

import (
	"context"
	"errors"
	"go-chat-app/models"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func CheckUserType(r *http.Request, role string) error {
	userType := r.Context().Value("user_type").(string)
	if userType != role {
		return errors.New("Unauthorized to access this resource")
	}
	return nil
}

func MatchUserTypeToUid(r *http.Request, userId string) error {
	userType := r.Context().Value("user_type").(string)
	uid := r.Context().Value("uid").(string)

	if userType == "USER" && uid != userId {
		return errors.New("Unauthorized to access this resource")
	}
	return CheckUserType(r, userType)
}


func GetUserTypeByID(userID string) (string, error) {
    var user models.User
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err := userCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
    if err != nil {
        return "", err
    }

    return *user.User_type, nil
}
