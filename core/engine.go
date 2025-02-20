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
	Db				   *Database
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
		Db:                 NewDatabase(),
		ProcessDefinitions: make(map[string]Definitions),
		ProcessModels:      make(map[string]Process),
		ProcessInstances: 	make(map[string]ProcessInstance),

	}
}

// Starts the engine
func (e *Engine) Start() {
	fmt.Println("Starting Engine...")
	fmt.Println("Initializing Database...")
	e.Db.InitializeDB()
	fmt.Println("Initializing Router...")
	e.InitRouter()
	fmt.Println("Loading ProcessModels...")
	e.LoadProcessModels()
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

	err := e.Db.SaveDefinitionToDB(definition)

	return err
}

func (e *Engine) AddProcessDefinition(definition *Definitions) {
	e.ProcessDefinitions[definition.XMLName.Local] = *definition

	for _, process := range definition.Processes {
		e.ProcessModels[process.ID] = process
	}

	e.Db.SaveDefinitionToDB(definition)
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

func (e *Engine) LoadProcessModels() {
	xmls, err := e.Db.LoadAllXMLs()
	
	if err != nil {
		log.Fatalf("Cannot get ProcessDefinitions from DB: %v", err)
	}

	for _, xmlString := range xmls {
		definition, err := ParseFromBpmnString(xmlString)
		if err != nil {
			log.Printf("Failed to parse process definition: %v", err)
			continue
		}

		for _, process := range definition.Processes {
			e.ProcessModels[process.ID] = process
		}
	}

	log.Println("Process models loaded successfully")
}

