package api

import (
	"encoding/json"
	"time"
)

type EngineInfo struct {
	Name string
	Port int
}

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
