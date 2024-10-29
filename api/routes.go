package api

import (
	"net/http"

	"github.com/luisthieme/GoMotion/pkg/middleware"
)


func RegisterRoutes() {
	http.HandleFunc("/user/profile", middleware.ApplyMiddlewares(HandleClientProfile))
	// http.HandleFunc("/user/settings", applyMiddlewares(HandleUserSettings)) // New endpoint
	// http.HandleFunc("/user/logout", applyMiddlewares(HandleUserLogout)) // Another new endpoint
}
