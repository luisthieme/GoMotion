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
	r.Mux.HandleFunc("/ws", r.Engine.EventManager.HandleConnections)

	r.Mux.HandleFunc("/", r.HandleBase)
	r.Mux.HandleFunc("GET /hello", r.HandleHello)
	r.Mux.HandleFunc("GET /go_motion/api/v1/info", r.HandleEngineInfo)
	r.Mux.HandleFunc("POST /go_motion/api/v1/start/{processModelId}", r.HandleStartProcessModel)
	r.Mux.HandleFunc("POST /go_motion/api/v1/process_definitions", r.HandleDeployProcessModel)
	r.Mux.HandleFunc("/go_motion/api/v1/process_instances", r.HandleProcessInstances)
	r.Mux.HandleFunc("GET /go_motion/api/v1/process_models", r.HandleProcessModels)
}

func (r *Router) HandleBase(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Go-Motion Workflow Engine</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					display: flex;
					justify-content: center;
					align-items: center;
					height: 100vh;
					margin: 0;
					background-color: #f0f0f0;
				}
				.container {
					text-align: center;
					padding: 2rem;
					background-color: white;
					border-radius: 8px;
					box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
				}
				h1 {
					color: #333;
				}
				p {
					color: #666;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Welcome to Go-Motion</h1>
				<p>The Engine is up and running!</p>
			</div>
		</body>
		</html>
	`

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func (r *Router) HandleEngineInfo(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := EngineInfo{
		Name: r.Engine.Name,
		Version: r.Engine.Version,
	}
	

	json.NewEncoder(w).Encode(response)
}

func (r *Router) HandleHello(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Received Hello!")
}

func (r *Router) HandleStartProcessModel(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	processModelId := req.PathValue("processModelId")

	r.Engine.StartProcess(processModelId)
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

	r.Engine.AddProcessDefinition(definitions)

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

func (r *Router) HandleProcessInstances(w http.ResponseWriter, req *http.Request) {
	// CORS-Header setzen
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// OPTIONS-Anfragen für CORS vorab genehmigen
	if req.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch req.Method {
	case http.MethodGet:
		rows, err := r.Engine.Db.Db.Query("SELECT id, process_model_name, state FROM process_instances")
		if err != nil {
			http.Error(w, "Failed to query process instances", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var processInstances []ProcessInstanceApiResponse
		for rows.Next() {
			var pi ProcessInstanceApiResponse
			err := rows.Scan(&pi.Id, &pi.ProcessModel, &pi.State)
			if err != nil {
				http.Error(w, "Failed to scan process instance", http.StatusInternalServerError)
				return
			}
			processInstances = append(processInstances, pi)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating through process instances", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(processInstances); err != nil {
			http.Error(w, "Failed to encode process instances as JSON", http.StatusInternalServerError)
			return
		}

	case http.MethodDelete:
		_, err := r.Engine.Db.Db.Exec("DELETE FROM process_instances")
		if err != nil {
			http.Error(w, "Failed to delete process instances", http.StatusInternalServerError)
			return
		}

		r.Engine.ProcessInstances = make(map[string]ProcessInstance)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "All process instances deleted successfully"}`))

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}


func (r *Router) HandleProcessModels(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Extract process models from the engine
	processModels := []map[string]string{}

	for _, process := range r.Engine.ProcessModels {
		processModels = append(processModels, map[string]string{
			"ID":      process.ID,
			"XMLName": process.XMLName.Local,
		})
	}

	// Encode the result as JSON and send the response
	if err := json.NewEncoder(w).Encode(processModels); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

