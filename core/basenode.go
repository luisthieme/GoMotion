package core

type BaseNode struct {
	BpmnType string
}

func (n *BaseNode) GetBpmnType() string {
	return n.BpmnType
}
