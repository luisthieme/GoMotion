package core

import "fmt"

type Engine struct {
	Name             string
	Url              string
	ProcessInstances map[string]ProcessInstance
}

// Start starts the engine
func (e Engine) Start() {
	fmt.Println("Starting Engine...")
}

// GetAllProcessInstances returns all process instances
func (e Engine) GetAllProcessInstances() map[string]ProcessInstance {
	return e.ProcessInstances
}

// GetProcessInstanceById retrieves a process instance by ID
func (e Engine) GetProcessInstanceById(id string) (ProcessInstance, error) {
	processInstance, exists := e.ProcessInstances[id]
	if !exists {
		return ProcessInstance{}, fmt.Errorf("process instance with ID %s not found", id)
	}
	return processInstance, nil
}

// StartProcessInstance starts the process instance with the specified ID
func (e Engine) StartProcessInstance(id string) (ProcessInstance, error) {
	processInstance, exists := e.ProcessInstances[id]
	if !exists {
		return ProcessInstance{}, fmt.Errorf("process instance with ID %s not found", id)
	}

	err := processInstance.Start();

	if err != nil {
		return ProcessInstance{}, err
	}

	return processInstance, nil
}

func (e Engine) AddProcessInstance(processInstance *ProcessInstance) error {
	e.ProcessInstances[processInstance.Id] = *processInstance

	return nil
}
