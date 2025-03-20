package core

import (
	"encoding/xml"
	"fmt"
	"os"
)

func ParseBpmnFromFile(filePath string) (*Definitions, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	var definitions Definitions
	err = decoder.Decode(&definitions)
	if err != nil {
		return nil, fmt.Errorf("failed to parse BPMN XML: %v", err)
	}

	return &definitions, nil
}

func ParseFromBpmnString(xmlString string) (*Definitions, error) {
	var definitions Definitions
	err := xml.Unmarshal([]byte(xmlString), &definitions)
	fmt.Printf("parsed bpmn definition: %s", definitions)
	if err != nil {
		return nil, fmt.Errorf("failed to parse xml: %w", err)
	}

	return &definitions, nil
}
