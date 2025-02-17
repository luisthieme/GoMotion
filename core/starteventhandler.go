package core

import "fmt"

type StartEventHandler struct {
	StartEvent *StartEvent
}

func NewStartEventHanler(startevent *StartEvent) *StartEventHandler {
	return &StartEventHandler{ StartEvent: startevent}
}

func (s *StartEventHandler) Execute() {
	fmt.Println("Executing Startevent")
}
