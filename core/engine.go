package core

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

// TODO: evaluate if sync.Map or sync.RWMutex or a normal ap should be used in the cases
type Engine struct {
	Name               string
	Port               string
	Version 	       string
	Router             Router
	EventManager       EventManager
	Db				   *Database
	ProcessDefinitions map[string]Definitions
	ProcessModels      map[string]ProcessModel
	ProcessInstances   map[string]ProcessInstance
	PendingTasks	   sync.Map
}


// Constructor
func NewEngine(name, port string) *Engine {
	return &Engine{
		Name:               name,
		Port:               port,
		Version: 			"0.0.1",
		Router: 			Router{},
		EventManager:       *NewEventManager(),
		Db:                 NewDatabase(),
		ProcessDefinitions: make(map[string]Definitions),
		ProcessModels:      make(map[string]ProcessModel),
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
	
	fmt.Println("Engine running on http://localhost:" + e.Port)
	log.Fatal(http.ListenAndServe(":" + e.Port, e.Router.Mux))
}

// Loads and parses the BPMN-File and saves it in the Engine and the connected DB
func (e *Engine) LoadAndAddProcessDefinition(filePath string) error {
	definition, error := ParseBpmnFromFile(filePath)

	if error != nil {
		return error
	}

	e.ProcessDefinitions[definition.XMLName.Local] = *definition

	for _, process := range definition.Processes {
		e.ProcessModels[process.ID] = ProcessModel{ Process: process, DefitionionId: definition.ID }
	}

	err := e.Db.SaveDefinitionToDB(definition)

	return err
}

// 
func (e *Engine) AddProcessDefinition(definition *Definitions) {
	e.ProcessDefinitions[definition.XMLName.Local] = *definition

	for _, process := range definition.Processes {
		e.ProcessModels[process.ID] = ProcessModel{ Process: process, DefitionionId: definition.ID }
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
			e.ProcessModels[process.ID] = ProcessModel{ Process: process, DefitionionId: definition.ID }
		}
	}

	log.Println("Process models loaded successfully")
}

// Register a pending task with a completion callback
func (e *Engine) RegisterPendingTask(taskID string, pendingTask PendingTask) {
	e.PendingTasks.Store(taskID, pendingTask)
	fmt.Println("Registered pending task:", taskID)
}

// Finish a task by calling the stored callback
func (e *Engine) CompletePendingTask(taskID string) bool {
	if taskValue, ok := e.PendingTasks.Load(taskID); ok {
		// Get the pending task from the map
		pendingTask, valid := taskValue.(PendingTask)
		if !valid {
			fmt.Println("Task found but invalid type:", taskID)
			return false
		}
		
		// Remove the task from the map before executing callback
		e.PendingTasks.Delete(taskID)
		
		// Execute the callback function if it exists
		if pendingTask.Callback != nil {
			pendingTask.Callback() // Execute the callback function
			fmt.Println("Completed task:", taskID, "name:", pendingTask.Name)
			return true
		} else {
			fmt.Println("Task found but callback is nil:", taskID)
			return false
		}
	} else {
		fmt.Println("Task not found:", taskID)
	}
	return false
}


