package core

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ProcessInstance struct {
	Id string
	ProcessModel ProcessModel
	CurrentElement string
	StartedAt time.Time
	FinishedAt time.Time
	State string
	Engine *Engine
}

func NewProcessInstance(processModel ProcessModel, engine *Engine) *ProcessInstance {
	return &ProcessInstance{
		Id:           uuid.NewString(),
		ProcessModel: processModel,
		Engine: engine,
	}
}

// Execute starts the process instance and executes events and tasks based on SequenceFlows.
func (p *ProcessInstance) Execute(token Token) error {
	fmt.Println("Executing process instance:", p.Id)

	if len(p.ProcessModel.StartEvents) < 1 {
		error := fmt.Errorf("no StartEvent in ProcessModel")

		return error
	}

	p.State = "running"
	p.StartedAt = time.Now()
	p.Engine.Db.SaveProcessInstanceToDB(p)

	p.Engine.EventManager.Broadcast(ProcessInstanceEvent{ Name: "started", Type: "processinstance", Id: p.Id, ProcessModelName: p.ProcessModel.Name, CurrentElement: p.CurrentElement, StartedAt: p.StartedAt, FinishedAt: p.FinishedAt})
	
	startEvent := p.ProcessModel.StartEvents[0]
	p.CurrentElement = startEvent.ID
	p.updateProcessInstanceState()
	startEventHandler := NewStartEventHanler(&startEvent, p)
	startEventHandler.Execute(token)
	p.CurrentElement = p.getNextElementFromFlow(startEvent.Outgoing)


	for {
		// Check if there is a new CurrentElement or if the ProcessInstance is ending
		if p.CurrentElement == "" {
			fmt.Println("ProcessInstance finished.")
			p.State = "finished"
			p.FinishedAt = time.Now()
			p.Engine.Db.PersistProcessInstance(p)
			p.Engine.EventManager.Broadcast(ProcessInstanceEvent{ Name: "finished", Type: "processinstance", Id: p.Id, ProcessModelName: p.ProcessModel.Name, CurrentElement: p.CurrentElement, StartedAt: p.StartedAt, FinishedAt: p.FinishedAt})
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
			endEventHandler := NewEndEventHandler(endEvent, p)
			endEventHandler.Execute(token)
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
			p.updateProcessInstanceState()
			taskHandler := NewTaskHandler(task, p)
			taskHandler.Execute(token)
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

func (p *ProcessInstance) updateProcessInstanceState() {
	p.Engine.Db.PersistProcessInstance(p)
	p.Engine.EventManager.Broadcast(ProcessInstanceEvent{ Name: "running", Type: "processinstance", Id: p.Id, ProcessModelName: p.ProcessModel.Name, CurrentElement: p.CurrentElement, StartedAt: p.StartedAt, FinishedAt: p.FinishedAt})
}

