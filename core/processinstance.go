package core

import (
	"fmt"

	"github.com/google/uuid"
)

type ProcessInstance struct {
	Id string
	ProcessModel Process
	CurrentElement string
	Engine *Engine
}

func NewProcessInstance(processModel Process, engine *Engine) *ProcessInstance {
	return &ProcessInstance{
		Id:           uuid.NewString(),
		ProcessModel: processModel,
		Engine: engine,
	}
}

// Execute starts the process instance and executes events and tasks based on SequenceFlows.
func (p *ProcessInstance) Execute() error {
	fmt.Println("Executing process instance:", p.Id)

	if len(p.ProcessModel.StartEvents) < 1 {
		error := fmt.Errorf("no StartEvent in ProcessModel")

		return error
	}

	
	startEvent := p.ProcessModel.StartEvents[0]
	p.CurrentElement = startEvent.ID
	startEventHandler := NewStartEventHanler(&startEvent, p)
	startEventHandler.Execute()
	p.CurrentElement = p.getNextElementFromFlow(startEvent.Outgoing)


	for {
		// Check if there is a new CurrentElement or if the ProcessInstance is ending
		if p.CurrentElement == "" {
			fmt.Println("ProcessInstance finished.")
			break
		}

		// Check if CurrentElement is an EndEvent and execute it
		var endEvent *EndEvent

		for _, e := range p.ProcessModel.EndEvents {
			if e.ID == p.CurrentElement {
				endEvent = &e
				break
			}
		}

		if endEvent != nil {
			endEventHandler := NewEndEventHandler(endEvent)
			endEventHandler.Execute()
			p.CurrentElement = ""
		}

		// Check if the CurrentElement is a Task and execute it
		var task *Task
		
		for _, t := range p.ProcessModel.Tasks {
			if t.ID == p.CurrentElement {
				task = &t
				break
			}
		}

		if task != nil {
			taskHandler := NewTaskHandler(task)
			taskHandler.Execute()
			p.CurrentElement = p.getNextElementFromFlow(task.Outgoing)
		}
		
	}
	

	return nil
}

func (p *ProcessInstance) getNextElementFromFlow(outgoing []string) string {
	for _, sequenceFlowID := range outgoing {
		for _, flow := range p.ProcessModel.SequenceFlows {
			if flow.ID == sequenceFlowID {
				return flow.TargetRef
			}
		}
	}
	return ""
}

