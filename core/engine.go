package core

import (
	"fmt"
	"log"
	"net/http"
)

type Engine struct {
	Name               string
	Url                string
	Version 	       string
	Router             Router
	EventManager       EventManager
	ProcessDefinitions map[string]Definitions
	ProcessModels      map[string]Process
	ProcessInstances   map[string]ProcessInstance
}

func NewEngine(name, url string) *Engine {
	return &Engine{
		Name:               name,
		Url:                url,
		Version: 			"0.0.1",
		Router: 			Router{},
		EventManager:       *NewEventManager(),
		ProcessDefinitions: make(map[string]Definitions),
		ProcessModels:      make(map[string]Process),
		ProcessInstances: 	make(map[string]ProcessInstance),

	}
}

// Starts the engine
func (e *Engine) Start() {
	fmt.Println("Starting Engine...")
	e.InitRouter()
	fmt.Println("Engine running on http://localhost:6969")
	log.Fatal(http.ListenAndServe(":6969", e.Router.Mux))
}

// Loads and parses the BPMN-File and saves it in the Engine
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

func (e *Engine) AddProcessDefinition(definition *Definitions) {
	e.ProcessDefinitions[definition.XMLName.Local] = *definition

	for _, process := range definition.Processes {
		e.ProcessModels[process.ID] = process
	}
}

func (e *Engine) InitRouter() {
	e.Router.Engine = e

	mux := http.NewServeMux()
	e.Router.Mux = mux

	e.Router.RegisterRoutes()
}

// Starts a ProcessInstance for a given ProcessModel
func (e *Engine) StartProcess(processModelId string) error {
	processModel := e.ProcessModels[processModelId]

	processInstance := NewProcessInstance(processModel, e)

	e.ProcessInstances[processInstance.Id] = *processInstance

	go processInstance.Execute()

	return nil
}
