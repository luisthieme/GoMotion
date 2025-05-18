package core

import (
	"fmt"

	"github.com/google/uuid"
)

type StartEventHandler struct {
	Id string
	StartEvent *StartEvent
	ProcessInstance *ProcessInstance
}

func NewStartEventHanler(startevent *StartEvent, processInstance *ProcessInstance) *StartEventHandler { 
	return &StartEventHandler{Id: uuid.NewString(), StartEvent: startevent, ProcessInstance: processInstance}
}

func (s *StartEventHandler) Execute(token Token) {
	s.ProcessInstance.Engine.EventManager.Broadcast(Event{ Name: "executing", Type: "startevent", Id: s.Id})
	fmt.Println("Executing Startevent")
	s.ProcessInstance.Engine.EventManager.Broadcast(Event{ Name: "finished", Type: "startevent", Id: s.Id})
}
