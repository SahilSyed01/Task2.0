package helpers

import (
	"errors"
	"net/http"
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
