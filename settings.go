package main

import (
	"encoding/json"
	"fmt"

	kubewarden "github.com/kubewarden/policy-sdk-go"
	kubewardenProtocol "github.com/kubewarden/policy-sdk-go/protocol"
)

// Settings defines the settings of the policy
type Settings struct {
	// Validations is the list of validations to be performed
	Validations []Validation `json:"validations"`
}

type Validation struct {
	// Message is the message to be shown to the user when the request is rejected
	Message string `json:"message"`
	// Expression is the CEL expression to be evaluated
	Expression string `json:"expression"`
}

func validateSettings(_ []byte) ([]byte, error) {
	return kubewarden.AcceptSettings()
}

func NewSettingsFromValidationReq(validationReq *kubewardenProtocol.ValidationRequest) (Settings, error) {
	settings := Settings{}

	if err := json.Unmarshal(validationReq.Settings, &settings); err != nil {
		return Settings{}, fmt.Errorf("cannot unmarshal settings %w", err)
	}
	return settings, nil
}
