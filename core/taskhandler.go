package core

import (
	"fmt"

	"github.com/google/uuid"
)

type TaskHandler struct {
	Id              string
	Task            *Task
	ProcessInstance *ProcessInstance
	Completed       chan bool
}

type PendingTask struct {
	Name              string
	ProcessModel      string
	ProcessInstanceId string
	Callback          func()
}

func NewTaskHandler(task *Task, processInstance *ProcessInstance) *TaskHandler {
	return &TaskHandler{
		Id:              uuid.NewString(),
		Task:            task,
		ProcessInstance: processInstance,
		Completed:       make(chan bool),
	}
}

// Execute a Task and wait for a completion signal via the API
func (t *TaskHandler) Execute() {
	t.ProcessInstance.Engine.EventManager.Broadcast(TaskEvent{Name: "executing", Type: "task", Id: t.Id, ElementName: t.Task.Name, ProcessInstanceId: t.ProcessInstance.Id})
	fmt.Println("Task is now waiting for completion...")

	t.ProcessInstance.Engine.RegisterPendingTask(t.Id, PendingTask{Name: t.Task.Name, ProcessModel: t.ProcessInstance.ProcessModel.Name, ProcessInstanceId: t.ProcessInstance.Id, Callback: func() {
		t.Completed <- true
	}} )

	t.ProcessInstance.Engine.EventManager.Broadcast((TaskEvent{Name: "pending", Type: "task", Id: t.Id, ElementName: t.Task.Name, ProcessInstanceId: t.ProcessInstance.Id}))

	<-t.Completed

	t.ProcessInstance.Engine.EventManager.Broadcast(TaskEvent{Name: "finished", Type: "task", Id: t.Id, ElementName: t.Task.Name, ProcessInstanceId: t.ProcessInstance.Id})
	fmt.Println("Task completed!")
}
