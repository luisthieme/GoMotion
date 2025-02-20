package core

type EngineInfo struct {
	Name string
	Version string
}

type ProcessInstanceApiResponse struct {
	Id            string `json:"id"`
	ProcessModel  string `json:"process_model"`
	State         string `json:"state"`
}
