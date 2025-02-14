package core

import (
	"fmt"

	"github.com/google/uuid"
)

type ProcessInstance struct {
	Id uuid.UUID
	ProcessModel Process
}

//Start starts the process instance
func (p ProcessInstance) Start() error {
	fmt.Println("Starting ProcessInstance...")

	p.Execute()
	
	return nil
}

//Execute shedules the execution of the next task
func (p ProcessInstance) Execute() error {
	fmt.Printf("Executing ProcessInstance...")

	startEvents := p.ProcessModel.GetStartEvents()

	startEventHandler := NewStartEventHandler(startEvents[0])

	startEventHandler.Execute()

	return nil
} 
