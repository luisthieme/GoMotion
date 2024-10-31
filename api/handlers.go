package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/luisthieme/GoMotion/internal"
)

// Endpoint: INFO
type EngineInfo struct {
	Name string
	Port int
}

// Only GET is allowed as a Method so we dont need seperate handlers
func HandleEngineInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	w.Header().Set("Content-Type", "application/json")

	response := EngineInfo{
		Name: internal.EngineName,
		Port: internal.Port,
	}
	

	json.NewEncoder(w).Encode(response)
}

// Endpoint: PROCESS_MODELS
type ProcessModel struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description,omitempty"`
	Version      int             `json:"version"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	IsExecutable bool            `json:"is_executable"`
	ProcessData  json.RawMessage `json:"process_data"`
}

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
	rows, err := internal.Db.Query("SELECT id, name, description, version, created_at, updated_at, is_executable, process_data FROM processmodels")
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var processModels []ProcessModel
	for rows.Next() {
		var pd ProcessModel
		err := rows.Scan(&pd.ID, &pd.Name, &pd.Description, &pd.Version, &pd.CreatedAt, &pd.UpdatedAt, &pd.IsExecutable, &pd.ProcessData)
		if err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		processModels = append(processModels, pd)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(processModels); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// POST a new ProcessModel
func PostProcessModels(w http.ResponseWriter, r *http.Request) {
	var pd ProcessModel

	if err := json.NewDecoder(r.Body).Decode(&pd); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	pd.CreatedAt = time.Now()
	pd.UpdatedAt = pd.CreatedAt
	pd.ID = uuid.New().String()

	query := `INSERT INTO processmodels (id, name, description, version, created_at, updated_at, is_executable, process_data) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := internal.Db.Exec(query, pd.ID, pd.Name, pd.Description, pd.Version, pd.CreatedAt, pd.UpdatedAt, pd.IsExecutable, pd.ProcessData)
	if err != nil {
		http.Error(w, "Failed to insert record", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(pd); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
