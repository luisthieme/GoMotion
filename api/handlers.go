package api

import (
	"encoding/json"
	"net/http"

	"github.com/luisthieme/GoMotion/internal"
	"github.com/luisthieme/GoMotion/pkg/middleware"
)

func HandleClientProfile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetClientProfile(w,r)
	case http.MethodPatch:
		UpdateClientProfile(w,r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}	
}

func GetClientProfile(w http.ResponseWriter, r *http.Request) {
	clientProfile := r.Context().Value(middleware.ClientProfileKey).(internal.ClientProfile)

	w.Header().Set("Content-Type", "application/json")

	response := internal.ClientProfile {
		Email: clientProfile.Email,
		Name: clientProfile.Name,
		Id: clientProfile.Id,
	}

	json.NewEncoder(w).Encode(response)
}

func UpdateClientProfile(w http.ResponseWriter, r *http.Request) {
	clientProfile := r.Context().Value(middleware.ClientProfileKey).(internal.ClientProfile)

	var payloadData internal.ClientProfile
	err := json.NewDecoder(r.Body).Decode(&payloadData)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if payloadData.Email != "" {
		clientProfile.Email = payloadData.Email
	}

	if payloadData.Name != "" {
		clientProfile.Name = payloadData.Name
	}

	internal.Database[clientProfile.Id] = clientProfile

	w.WriteHeader(http.StatusOK)
}
