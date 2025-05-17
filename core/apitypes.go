package core

import "time"

type EngineInfo struct {
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	StartedAt time.Time `json:"started_at"`
}

type ProcessInstanceApiResponse struct {
	Id               string    `json:"id"`
	ProcessModelName string    `json:"process_model_name"`
	CurrentElement   string    `json:"current_element"`
	StartedAt        time.Time `json:"started_at"`
	FinishedAt       *time.Time `json:"finished_at"`
	State            string    `json:"state"`
}
