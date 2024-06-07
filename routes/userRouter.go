package routes

import (
	"go-chat-app/controllers"
	"go-chat-app/middleware"
	"net/http"
)

func UserRoutes() {
	// Only protect the routes that should require authentication
	http.Handle("/users", middleware.Authenticate(http.HandlerFunc(controllers.GetUsers)))
	http.Handle("/users/", middleware.Authenticate(http.HandlerFunc(controllers.GetUser)))
}
