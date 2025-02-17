package core

import "fmt"

type EndEventHandler struct {
	EndEvent *EndEvent
}

func NewEndEventHandler(endEvent *EndEvent) *EndEventHandler{
	return &EndEventHandler{ EndEvent: endEvent}
} 

func (e *EndEventHandler) Execute() {
	fmt.Println("Executing EndEvent")
}
