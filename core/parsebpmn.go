package core

import (
	"encoding/xml"
	"fmt"
	"strings"
)

func ParseBpmnFromString(xmlString string, definition *BPMNDefinitions) error{

	reader := strings.NewReader(xmlString)
	decoder := xml.NewDecoder(reader)

	if err := decoder.Decode(&definition); err != nil {
		fmt.Println("Error decoding XML: ", err)
		return err
	}

	return nil
}
