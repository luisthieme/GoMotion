package api

import (
	"net/http"
)

// basePath := "/go_motion/api/v1/"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /go_motion/api/v1/info", HandleEngineInfo)
	mux.HandleFunc("/go_motion/api/v1/process_models", HandleProcessModels)
	mux.HandleFunc("/go_motion/api/v1/process_models/{processModelId}", HandleProcessModel)
}
