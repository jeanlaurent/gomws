package main

import (
	"encoding/xml"
	"errors"
	"strings"
)

//StripErrorMessage of MWS Error response
func stripErrorMessage(someXML string) string {
	errorMessage, error := stripSingleString(someXML, "Message")
	if error != nil {
		return "Unknown Error"
	}
	return errorMessage
}

func stripSingleString(someXML string, stringID string) (string, error) {
	decoder := xml.NewDecoder(strings.NewReader(someXML))
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch tokenType := token.(type) {
		case xml.StartElement:
			if tokenType.Name.Local == stringID {
				var message string
				decoder.DecodeElement(&message, &tokenType)
				return message, nil
			}
		}
	}
	return "", errors.New("not found")
}
