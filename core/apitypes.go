package core

type EngineInfo struct {
	Name    string
	Version string
}

type ProcessInstanceApiResponse struct {
	Id               string `json:"id"`
	ProcessModelName string `json:"process_model_name"`
	State            string `json:"state"`
}
