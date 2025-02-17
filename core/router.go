package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Router struct {
	Engine *Engine
	Mux *http.ServeMux
}

func (r *Router) RegisterRoutes() {
	r.Mux.HandleFunc("GET /hello", r.HandleHello)
	r.Mux.HandleFunc("GET /go_motion/api/v1/info", r.HandleEngineInfo)
	r.Mux.HandleFunc("POST /go_motion/api/v1/start/{processModelId}", r.HandleStartProcessModel)
	r.Mux.HandleFunc("POST /go_motion/api/v1/process_definitions", r.HandleDeployProcessModel)
}

func (router *Router) HandleEngineInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := EngineInfo{
		Name: router.Engine.Name,
		Version: router.Engine.Version,
	}
	

	json.NewEncoder(w).Encode(response)
}

func (router *Router) HandleHello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Hello!")
}

func (router *Router) HandleStartProcessModel(w http.ResponseWriter, r *http.Request) {
	processModelId := r.PathValue("processModelId")

	router.Engine.StartProcess(processModelId)
}

func (r *Router) HandleDeployProcessModel(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Failed to read request body.", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	// JSON-Daten parsen
	var payload struct {
		XML string `json:"xml"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Error while parsing JSON-Request", http.StatusBadRequest)
		return
	}

	// BPMN-XML parsen
	definitions, err := ParseFromBpmnString(payload.XML)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while parsing XML: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("A new BPMN was deployed.")

	// Erfolgreiche Antwort zurückgeben
	response := struct {
		Message     string       `json:"message"`
		Processes   []Process    `json:"processes"`
	}{
		Message:   "Successfully deployed BPMN.",
		Processes: definitions.Processes,
	}

	json.NewEncoder(w).Encode(response)
}
