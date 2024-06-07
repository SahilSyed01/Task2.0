package middleware

import (
	"context"
	"net/http"

	"go-chat-app/helpers"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientToken := r.Header.Get("Authorization")
		if clientToken == "" {
			http.Error(w, "No Authorization header provided", http.StatusUnauthorized)
			return
		}

		claims, err := helpers.ValidateToken(clientToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "first_name", claims.First_name)
		ctx = context.WithValue(ctx, "uid", claims.Uid)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
