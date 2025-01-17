package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/luisthieme/GoMotion/core"
	"github.com/luisthieme/GoMotion/internal"
)

// Endpoint: INFO

// Only GET is allowed as a Method so we dont need seperate handlers
func HandleEngineInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := EngineInfo{
		Name: internal.EngineName,
		Port: internal.Port,
	}
	

	json.NewEncoder(w).Encode(response)
}

// Endpoint: PROCESS_MODELS

// Handlers for ProcessModels
func HandleProcessModels(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetProcessModels(w,r)
	case http.MethodPost:
		PostProcessModels(w,r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET all ProcessModels
func GetProcessModels(w http.ResponseWriter, r *http.Request) {
	// rows, err := internal.Db.Query("SELECT * FROM processmodels")
	// if err != nil {
	// 	http.Error(w, "Failed to query database", http.StatusInternalServerError)
	// 	return
	// }
	// defer rows.Close()

	// var processModels []ProcessModel
	// for rows.Next() {
	// 	var pd ProcessModel
	// 	err := rows.Scan(&pd.ID, &pd.Name, &pd.Description, &pd.Version, &pd.CreatedAt, &pd.UpdatedAt, &pd.IsExecutable, &pd.XML)
	// 	if err != nil {
	// 		http.Error(w, "Failed to scan row", http.StatusInternalServerError)
	// 		return
	// 	}
	// 	processModels = append(processModels, pd)
	// }

	// w.Header().Set("Content-Type", "application/json")
	// if err := json.NewEncoder(w).Encode(processModels); err != nil {
	// 	http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	// }
}

// POST a new ProcessModel
func PostProcessModels(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST request for process_models received.")
	type RequestPayload struct {
		XML string `json:"xml"`
	}

	var reqPayload RequestPayload

	if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var bpmnDefinition core.BPMNDefinitions

	if err := core.ParseBpmnFromString(reqPayload.XML, &bpmnDefinition); err != nil {
		http.Error(w, "Failed to parse XML", http.StatusInternalServerError)
		return
	}

	var pd ProcessModel

	pd.CreatedAt = time.Now()
	pd.UpdatedAt = pd.CreatedAt
	pd.ID = bpmnDefinition.ID
	pd.Name = bpmnDefinition.XMLName.Local
	pd.XML = reqPayload.XML
	pd.IsExecutable = true

	query := `INSERT INTO processmodels (id, name, created_at, updated_at, is_executable, xml) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := internal.Db.Exec(query, pd.ID, pd.Name, pd.CreatedAt, pd.UpdatedAt, pd.IsExecutable, pd.XML)
	if err != nil {
		http.Error(w, "Failed to insert record", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(pd); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
	fmt.Printf("process model with id %s deployed.", pd.ID)
}

// Endpoint: PROCESS_MODEL

// Handler for ProcessModel
func HandleProcessModel(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetProcessModel(w,r)
	case http.MethodDelete:
		DeleteProcessModel(w,r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetProcessModel(w http.ResponseWriter, r *http.Request) {
	// processModelId := r.PathValue("processModelId")

	// var processModel ProcessModel
	// err := internal.Db.QueryRow(
	// 	"SELECT * FROM processmodels WHERE id=$1",
	// 	processModelId,
	// ).Scan(
	// 	&processModel.ID,
	// 	&processModel.Name,
	// 	&processModel.Description,
	// 	&processModel.Version,
	// 	&processModel.CreatedAt,
	// 	&processModel.UpdatedAt,
	// 	&processModel.IsExecutable,
	// 	&processModel.XML,
	// )
	
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		http.Error(w, "No process model found with the given ID", http.StatusNotFound)
	// 		return
	// 	}
	// 	http.Error(w, "Failed to query database", http.StatusInternalServerError)
	// 	return
	// }


	// w.Header().Set("Content-Type", "application/json")
	// if err := json.NewEncoder(w).Encode(processModel); err != nil {
	// 	http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	// }
}

func DeleteProcessModel(w http.ResponseWriter, r *http.Request) {
	processModelId := r.PathValue("processModelId")

	result, err := internal.Db.Exec("DELETE FROM processmodels WHERE id=$1", processModelId)

	if err != nil {
		http.Error(w, "Failed to execute delete query", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No process model found with the given ID", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent) 
}

// POST to start a ProcessModel

func StartProcessModel(w http.ResponseWriter, r *http.Request) {
	
}
