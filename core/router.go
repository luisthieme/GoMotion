package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Router struct {
	Engine *Engine
	Mux *http.ServeMux
}

func (r *Router) RegisterRoutes() {
	basePath := "/go_motion/api/v1"
	r.Mux.HandleFunc("/ws", r.Engine.EventManager.HandleConnections)
	r.Mux.HandleFunc("/", r.HandleBase)

	r.Mux.HandleFunc(basePath + "/info", r.HandleEngineInfo)

	r.Mux.HandleFunc(basePath + "/process_definitions", r.HandleDeployProcessModel)
	r.Mux.HandleFunc(basePath + "/process_models", r.HandleProcessModels)

	r.Mux.HandleFunc(basePath + "/start/{processModelId}", r.HandleStartProcessModel)
	
	r.Mux.HandleFunc(basePath + "/process_instances", r.HandleProcessInstances)

	r.Mux.HandleFunc(basePath + "/tasks/{taskId}/complete", r.HandleTaskCompletion)

	r.Mux.HandleFunc(basePath + "/tasks", r.HandleTasks)
}

func (r *Router) HandleBase(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	filePath := filepath.Join("html", "index.html") 

	html, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Could not read HTML file", http.StatusInternalServerError)
		fmt.Println("Error reading file:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(html)
}

func (r *Router) HandleEngineInfo(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	response := EngineInfo{
		Name: r.Engine.Name,
		Version: r.Engine.Version,
		StartedAt: r.Engine.StartedAt,
	}
	
	json.NewEncoder(w).Encode(response)
}

func (r *Router) HandleStartProcessModel(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	processModelId := req.PathValue("processModelId")

	r.Engine.StartProcess(processModelId)
}

func (r *Router) HandleDeployProcessModel(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
    // Handle preflight OPTIONS request
    if req.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

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
			err := rows.Scan(&pi.Id, &pi.ProcessModelName, &pi.State, &pi.StartedAt, &pi.FinishedAt)
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
			"name": process.Name,
			"id": process.ID,
			"definition_id": process.DefitionionId,
		})
	} 

	// Encode the result as JSON and send the response
	if err := json.NewEncoder(w).Encode(processModels); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (r *Router) HandleTasks(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// Check that we're handling a GET request
	if req.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Initialize an empty slice for task details
	taskDetails := []map[string]string{}
	
	r.Engine.PendingTasks.Range(func(key, value interface{}) bool {
		// Cast the value to a PendingTask
		t, ok := value.(PendingTask)
		if !ok {
			return true // Skip this item if casting fails
		}
		
		// Get the ID from the key
		id, ok := key.(string)
		if !ok {
			return true // Skip this item if key is not a string
		}
		
		// Create a task detail object in the requested format
		taskDetail := map[string]string{
			"name":              "pending",
			"type":              "task",
			"id":                id,
			"element_name":       t.Name,
			"process_instance_id": t.ProcessInstanceId,
		}
		
		taskDetails = append(taskDetails, taskDetail)
		return true
	})
	
	// This will return "[]" if taskDetails is empty
	err := json.NewEncoder(w).Encode(taskDetails)
	if err != nil {
		http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}

func (r *Router) HandleTaskCompletion(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	taskId := req.PathValue("taskId")
	success := r.Engine.CompletePendingTask(taskId)

	if success {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Task completed successfully"}`))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Task not found or already completed"}`))
	}
}


