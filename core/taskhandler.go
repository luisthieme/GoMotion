package core

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type TaskHandler struct {
	Id string
	Task *Task
	ProcessInstance *ProcessInstance
}

func NewTaskHandler(task *Task, processInstance *ProcessInstance) *TaskHandler {
	return &TaskHandler{ Id: uuid.NewString(), Task: task, ProcessInstance: processInstance}
}

func (t *TaskHandler) Execute() {
	t.ProcessInstance.Engine.EventManager.Broadcast(Event{ Name: "executing", Type: "task", Id: t.Id})
	fmt.Println("Executing Task")
	time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)
	t.ProcessInstance.Engine.EventManager.Broadcast(Event{ Name: "finished", Type: "task", Id: t.Id})
}
