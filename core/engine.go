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
	ProcessDefinitions map[string]Definitions
	ProcessModels      map[string]Process
	ProcessInstances   map[string]ProcessInstance
}

func NewEngine(name, url string) *Engine {
	return &Engine{
		Name:               name,
		Url:                url,
		Version: 			"0.0.1",
		ProcessDefinitions: make(map[string]Definitions),
		ProcessModels:      make(map[string]Process),
		ProcessInstances: 	make(map[string]ProcessInstance),
		Router: 			Router{},
	}
}

// Starts the engine
func (e *Engine) Start() {
	fmt.Println("Starting Engine...")
	e.LoadAndAddProcessDefinition("/Users/guestuser/5minds/PrivateStuff/GoMotion/processes/diagram_2.bpmn")
	e.InitRouter()
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

func (e *Engine) InitRouter() {
	e.Router.Engine = e

	mux := http.NewServeMux()
	e.Router.Mux = mux

	e.Router.RegisterRoutes()

	log.Fatal(http.ListenAndServe(":6969", mux))
}

// Starts a ProcessInstance for a given ProcessModel
func (e *Engine) StartProcess(processModelId string) error {
	processModel := e.ProcessModels[processModelId]

	processInstance := NewProcessInstance(processModel)

	e.ProcessInstances[processInstance.Id] = *processInstance

	go processInstance.Execute()

	return nil
}
