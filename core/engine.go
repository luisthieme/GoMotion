package core

import (
	"fmt"
)

type Engine struct {
	Name             string
	Url              string
	ProcessDefinitions map[string]Definitions
	ProcessModels map[string]Process
	ProcessInstances map[string]ProcessInstance
}

func NewEngine(name, url string) *Engine {
	return &Engine{
		Name:               name,
		Url:                url,
		ProcessDefinitions: make(map[string]Definitions),
		ProcessModels:      make(map[string]Process),
		ProcessInstances: 	make(map[string]ProcessInstance),
	}
}

// Start starts the engine
func (e *Engine) Start() {
	fmt.Println("Starting Engine...")
	e.LoadAndAddProcessDefinition("/Users/guestuser/5minds/PrivateStuff/GoMotion/processes/diagram.bpmn")
	e.StartProcess("Process_108m3pl")
}

func (e *Engine) LoadAndAddProcessDefinition(filePath string) error {
	definition, error := ParseBpmnFromFile(filePath)

	if error != nil {
		return error
	}

	e.ProcessDefinitions[definition.XMLName.Local] = *definition

	for _, process := range definition.Processes {
		e.ProcessModels[process.ID] = process
	}

	return nil
}

func (e *Engine) StartProcess(processModelId string) error {
	processModel := e.ProcessModels[processModelId]

	processInstance := NewProcessInstance(processModel)

	e.ProcessInstances[processInstance.Id] = *processInstance

	processInstance.Execute()

	return nil
}
