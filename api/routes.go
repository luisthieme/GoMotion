package api

import (
	"net/http"
)


func RegisterRoutes() {
	// http.HandleFunc("/user/profile", middleware.ApplyMiddlewares(HandleClientProfile))
	http.HandleFunc("go_motion/api/v1/info", HandleEngineInfo)
	http.HandleFunc("go_motion/api/v1/process_models", HandleProcessModels)
}
