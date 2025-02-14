package core

type StartEvent struct {
	BaseNode
	Name string
	Outgoing []string
}

func NewStartEvent(name string, outgoing []string) *StartEvent {
	return &StartEvent{
		BaseNode: BaseNode{ BpmnType: "start_event"},
		Name:     name,
		Outgoing: outgoing,
	}
}
