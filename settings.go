package main

import (
	"encoding/json"
	"log"
)

// Settings defines the settings of the policy
type Settings struct {
	// Message is the message to be shown to the user when the request is rejected
	Message string `json:"message"`
	// Expression is the CEL expression to be evaluated
	Expression string `json:"expression"`
}

func validateSettings(_ []byte) []byte {
	response := AcceptSettings()

	responseBytes, err := json.Marshal(&response)
	if err != nil {
		log.Fatalf("cannot marshal validation response: %v", err)
	}
	return responseBytes
}
