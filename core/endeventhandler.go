package core

import (
	"fmt"

	"github.com/google/uuid"
)

type EndEventHandler struct {
	Id string
	EndEvent *EndEvent
	ProcessInstance *ProcessInstance
}

func NewEndEventHandler(endEvent *EndEvent, processInstance *ProcessInstance) *EndEventHandler{
	return &EndEventHandler{ Id: uuid.NewString(), EndEvent: endEvent, ProcessInstance: processInstance}
} 

func (e *EndEventHandler) Execute(token Token) {
	e.ProcessInstance.Engine.EventManager.Broadcast(Event{ Name: "executing", Type: "endevent", Id: e.Id})
	fmt.Println("Executing EndEvent")
	e.ProcessInstance.Engine.EventManager.Broadcast(Event{ Name: "finished", Type: "endevent", Id: e.Id})
}
