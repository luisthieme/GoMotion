package api

import (
	"time"
)

type EngineInfo struct {
	Name string
	Port int
}

type ProcessModel struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	IsExecutable bool            `json:"is_executable"`
	XML  		 string          `json:"xml"`
}


