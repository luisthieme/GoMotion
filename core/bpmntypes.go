package core

import "encoding/xml"

// Definitions repräsentiert das Wurzelelement der BPMN-Datei
type Definitions struct {
	XMLName 	xml.Name  `xml:"definitions"`
	Processes   []Process `xml:"process"`
}

// Process enthält alle Flow-Elemente des BPMN-Modells
type Process struct {
	XMLName       xml.Name       `xml:"process"`
	ID            string         `xml:"id,attr"`
	StartEvents   []StartEvent   `xml:"startEvent"`
	Tasks         []Task         `xml:"task"`
	EndEvents     []EndEvent     `xml:"endEvent"`
	SequenceFlows []SequenceFlow `xml:"sequenceFlow"`
}

// StartEvent repräsentiert ein Start-Event
type StartEvent struct {
	XMLName  xml.Name   `xml:"startEvent"`
	ID       string     `xml:"id,attr"`
	Outgoing []string   `xml:"outgoing"`
}

// Task repräsentiert eine BPMN-Task
type Task struct {
	XMLName  xml.Name   `xml:"task"`
	ID       string     `xml:"id,attr"`
	Incoming []string   `xml:"incoming"`
	Outgoing []string   `xml:"outgoing"`
}

// EndEvent repräsentiert ein End-Event
type EndEvent struct {
	XMLName  xml.Name   `xml:"endEvent"`
	ID       string     `xml:"id,attr"`
	Incoming []string   `xml:"incoming"`
}

// SequenceFlow definiert die Verbindung zwischen BPMN-Elementen
type SequenceFlow struct {
	XMLName   xml.Name `xml:"sequenceFlow"`
	ID        string   `xml:"id,attr"`
	SourceRef string   `xml:"sourceRef,attr"`
	TargetRef string   `xml:"targetRef,attr"`
}
