package core

type Process struct {
	ID                   string               `json:"id"`
	Name                 *string              `json:"name,omitempty"`
	Version              *string              `json:"version,omitempty"`
	IsExecutable         bool                 `json:"isExecutable"`
	// LaneSet              *LaneSet             `json:"laneSet,omitempty"`
	FlowNodes            []FlowNode           `json:"flowNodes"`
	// SequenceFlows        []SequenceFlow       `json:"sequenceFlows"`
	// Associations         []Association        `json:"associations"`
	// DataStoreReferences  []DataStoreReference `json:"dataStoreReferences"`
	// DataObjectReferences []DataObjectReference `json:"dataObjectReferences"`
	// DataObjects          []DataObject         `json:"dataObjects"`
	// Documentation        *string              `json:"documentation,omitempty"`
	// ExtensionElements    *ExtensionElements   `json:"extensionElements,omitempty"`
	// EmbeddedSubProcesses []SubProcess         `json:"embeddedSubProcesses,omitempty"`
}

// GetStartEvents returns all StartEvents contained in this process.
func (p *Process) GetStartEvents() []*StartEvent {
	var startEvents []*StartEvent
	for _, node := range p.FlowNodes {
		if node.GetBpmnType() == "start_event" {
			startEvents = append(startEvents, node.(*StartEvent))
		}
	}
	return startEvents
}

// // GetEndEvents returns all EndEvents contained in this process.
// func (p *Process) GetEndEvents(includeEmbeddedSubProcess bool) []EndEvent {
// 	flowNodes := p.getRelevantFlowNodes(includeEmbeddedSubProcess)
// 	var endEvents []EndEvent
// 	for _, node := range flowNodes {
// 		if node.BpmnType == BpmnTypeEndEvent {
// 			endEvents = append(endEvents, node.(EndEvent))
// 		}
// 	}
// 	return endEvents
// }

// // GetFlowNodesByType returns all FlowNodes that match the given type.
// func (p *Process) GetFlowNodesByType(bpmnType BpmnType, includeEmbeddedSubProcess bool) []FlowNode {
// 	flowNodes := p.getRelevantFlowNodes(includeEmbeddedSubProcess)
// 	var matchingNodes []FlowNode
// 	for _, node := range flowNodes {
// 		if node.BpmnType == bpmnType {
// 			matchingNodes = append(matchingNodes, node)
// 		}
// 	}
// 	return matchingNodes
// }

// // GetAllEmbeddedSubProcesses returns all SubProcess activities in this process.
// func (p *Process) GetAllEmbeddedSubProcesses() []SubProcess {
// 	if len(p.EmbeddedSubProcesses) == 0 {
// 		var findSubSubProcesses func(subProcess SubProcess) []SubProcess
// 		findSubSubProcesses = func(subProcess SubProcess) []SubProcess {
// 			var subSubProcesses []SubProcess
// 			for _, node := range subProcess.FlowNodes {
// 				if node.BpmnType == BpmnTypeSubProcess {
// 					subSubProcesses = append(subSubProcesses, node.(SubProcess))
// 					subSubProcesses = append(subSubProcesses, findSubSubProcesses(node.(SubProcess))...)
// 				}
// 			}
// 			return append([]SubProcess{subProcess}, subSubProcesses...)
// 		}

// 		var embeddedSubProcesses []SubProcess
// 		for _, node := range p.FlowNodes {
// 			if node.BpmnType == BpmnTypeSubProcess {
// 				embeddedSubProcesses = append(embeddedSubProcesses, findSubSubProcesses(node.(SubProcess))...)
// 			}
// 		}
// 		p.EmbeddedSubProcesses = embeddedSubProcesses
// 	}
// 	return p.EmbeddedSubProcesses
// }

// // GetAllFlowNodes returns all FlowNodes of the model and all embedded models.
// func (p *Process) GetAllFlowNodes() []FlowNode {
// 	var allFlowNodes []FlowNode
// 	allFlowNodes = append(allFlowNodes, p.FlowNodes...)
// 	for _, subProcess := range p.GetAllEmbeddedSubProcesses() {
// 		allFlowNodes = append(allFlowNodes, subProcess.FlowNodes...)
// 	}
// 	return allFlowNodes
// }

// // GetAllSequenceFlows returns all SequenceFlows of the model and all embedded models.
// func (p *Process) GetAllSequenceFlows() []SequenceFlow {
// 	var allSequenceFlows []SequenceFlow
// 	allSequenceFlows = append(allSequenceFlows, p.SequenceFlows...)
// 	for _, subProcess := range p.GetAllEmbeddedSubProcesses() {
// 		allSequenceFlows = append(allSequenceFlows, subProcess.SequenceFlows...)
// 	}
// 	return allSequenceFlows
// }

// // GetAllDataObjectReferences returns all DataObjectReferences of the model and all embedded models.
// func (p *Process) GetAllDataObjectReferences() []DataObjectReference {
// 	var allDataObjectReferences []DataObjectReference
// 	allDataObjectReferences = append(allDataObjectReferences, p.DataObjectReferences...)
// 	for _, subProcess := range p.GetAllEmbeddedSubProcesses() {
// 		allDataObjectReferences = append(allDataObjectReferences, subProcess.DataObjectReferences...)
// 	}
// 	return allDataObjectReferences
// }

// // getRelevantFlowNodes returns either only this process's FlowNodes or all FlowNodes including embedded subprocesses.
// func (p *Process) getRelevantFlowNodes(includeEmbeddedSubProcess bool) []FlowNode {
// 	if includeEmbeddedSubProcess {
// 		return p.GetAllFlowNodes()
// 	}
// 	return p.FlowNodes
// }
