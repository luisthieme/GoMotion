package core

import "fmt"

type TaskHandler struct {
	Task *Task
}

func NewTaskHandler(task *Task) *TaskHandler {
	return &TaskHandler{ Task: task }
}

func (t *TaskHandler) Execute() {
	fmt.Println("Executing Task")
}
