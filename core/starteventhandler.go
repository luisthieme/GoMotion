package core

import (
	"fmt"

	"github.com/google/uuid"
)

type StartEventHandler struct {
	ID uuid.UUID
	*StartEvent
}


func NewStartEventHandler(startEvent *StartEvent) *StartEventHandler {
	return &StartEventHandler{
		ID: uuid.New(),
		StartEvent: startEvent,
		
	}
}

func (s *StartEventHandler) Execute() {
	fmt.Printf("ProcessInstance started with StartEvent: %s", s.Name)
	// fuer alle outgoing flows den richtigen handler erstellen

	// die execute function von allen oben erstellten handlern aufrufen handlern aufrufen
}

