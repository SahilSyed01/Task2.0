// middleware/authMiddleware.go

package middleware

import (
    // "context"
    "net/http"
    "go-chat-app/helpers"
)

func Authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        adminID := r.Header.Get("admin_id")
        if adminID == "" {
            http.Error(w, "No admin_id provided", http.StatusUnauthorized)
            return
        }

        userType, err := helpers.GetUserTypeByID(adminID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }

        if userType != "ADMIN" {
            http.Error(w, "Only admin has this privilege", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}
