package core

import (
	"encoding/xml"
)

type BPMNDefinitions struct {
	XMLName        xml.Name          `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL definitions"`
	TargetNamespace string           `xml:"targetNamespace,attr"`
	ID             string            `xml:"id,attr"`
	Collaboration  BPMNCollaboration `xml:"collaboration"`
	Process        BPMNProcess       `xml:"process"`
	Diagram        BPMNDiagram       `xml:"http://www.omg.org/spec/BPMN/20100524/DI BPMNDiagram"`
}

type BPMNCollaboration struct {
	ID          string           `xml:"id,attr"`
	Name        string           `xml:"name,attr"`
	Participants []BPMNParticipant `xml:"participant"`
}

type BPMNParticipant struct {
	ID         string `xml:"id,attr"`
	Name       string `xml:"name,attr"`
	ProcessRef string `xml:"processRef,attr"`
}

type BPMNProcess struct {
	ID           string          `xml:"id,attr"`
	Name         string          `xml:"name,attr"`
	IsExecutable bool            `xml:"isExecutable,attr"`
	LaneSet      BPMNLaneSet     `xml:"laneSet"`
	StartEvents  []BPMNStartEvent `xml:"startEvent"`
	EndEvents    []BPMNEndEvent   `xml:"endEvent"`
	ServiceTasks []BPMNServiceTask `xml:"serviceTask"`
	SequenceFlows []BPMNSequenceFlow `xml:"sequenceFlow"`
}

type BPMNLaneSet struct {
	Lanes []BPMNLane `xml:"lane"`
}

type BPMNLane struct {
	ID           string   `xml:"id,attr"`
	Name         string   `xml:"name,attr"`
	FlowNodeRefs []string `xml:"flowNodeRef"`
}

type BPMNStartEvent struct {
	ID       string   `xml:"id,attr"`
	Name     string   `xml:"name,attr"`
	Outgoing []string `xml:"outgoing"`
}

type BPMNEndEvent struct {
	ID       string   `xml:"id,attr"`
	Incoming []string `xml:"incoming"`
}

type BPMNServiceTask struct {
	ID       string   `xml:"id,attr"`
	Incoming []string `xml:"incoming"`
	Outgoing []string `xml:"outgoing"`
}

type BPMNSequenceFlow struct {
	ID        string `xml:"id,attr"`
	SourceRef string `xml:"sourceRef,attr"`
	TargetRef string `xml:"targetRef,attr"`
}

type BPMNDiagram struct {
	ID    string      `xml:"id,attr"`
	Plane BPMNPlane   `xml:"BPMNPlane"`
}

type BPMNPlane struct {
	ID          string         `xml:"id,attr"`
	BPMNElement string         `xml:"bpmnElement,attr"`
	Shapes      []BPMNShape    `xml:"BPMNShape"`
	Edges       []BPMNEdge     `xml:"BPMNEdge"`
}

type BPMNShape struct {
	ID          string  `xml:"id,attr"`
	BPMNElement string  `xml:"bpmnElement,attr"`
	IsHorizontal bool    `xml:"isHorizontal,attr,omitempty"`
	Bounds      DCBounds `xml:"http://www.omg.org/spec/DD/20100524/DC Bounds"`
	Label       *BPMNLabel `xml:"BPMNLabel,omitempty"`
}

type DCBounds struct {
	X      float64 `xml:"x,attr"`
	Y      float64 `xml:"y,attr"`
	Width  float64 `xml:"width,attr"`
	Height float64 `xml:"height,attr"`
}

type BPMNLabel struct {
	Bounds DCBounds `xml:"http://www.omg.org/spec/DD/20100524/DC Bounds"`
}

type BPMNEdge struct {
	ID          string    `xml:"id,attr"`
	BPMNElement string    `xml:"bpmnElement,attr"`
	Waypoints   []Waypoint `xml:"http://www.omg.org/spec/DD/20100524/DI waypoint"`
}

type Waypoint struct {
	X float64 `xml:"x,attr"`
	Y float64 `xml:"y,attr"`
}
