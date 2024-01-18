package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/ext"
	k8s "github.com/kubewarden/go-wasi-policy-template/cel/library"
	kubewarden "github.com/kubewarden/policy-sdk-go"
	kubewardenProtocol "github.com/kubewarden/policy-sdk-go/protocol"
)

func validate(payload []byte) ([]byte, error) {
	validationRequest := kubewardenProtocol.ValidationRequest{}

	err := json.Unmarshal(payload, &validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(fmt.Sprintf("Error deserializing validation request: %v", err)),
			kubewarden.Code(400))
	}

	env, err := cel.NewEnv(
		cel.EagerlyValidateDeclarations(true),
		cel.DefaultUTCTimeZone(true),
		ext.Strings(ext.StringsVersion(2)),
		cel.CrossTypeNumericComparisons(true),
		cel.OptionalTypes(),
		k8s.URLs(),
		k8s.Regex(),
		k8s.Lists(),
		cel.Variable("self", cel.DynType),
		cel.Variable("oldSelf", cel.DynType),
	)
	if err != nil {
		log.Fatalf("failed to create CEL env: %v", err)
	}

	settings, err := NewSettingsFromValidationReq(&validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(fmt.Sprintf("Error serializing RawMessage: %v", err)),
			kubewarden.Code(400))
	}

	for _, validation := range settings.Validations {
		ast, issues := env.Compile(validation.Expression)
		if issues != nil {
			log.Fatalf("failed to compile the CEL expression: %s", issues.String())
		}

		prog, err := env.Program(ast, cel.EvalOptions(cel.OptOptimize, cel.OptTrackCost))
		if err != nil {
			log.Fatalf("failed to instantiate CEL program: %v", err)
		}
		val, _, err := prog.Eval(map[string]interface{}{
			"self":    validationRequest.Request.Object,
			"oldSelf": validationRequest.Request.OldObject,
		})
		if err != nil {
			log.Fatalf("failed to evaluate: %v", err)
		}

		if val == types.False {
			return kubewarden.RejectRequest(kubewarden.Message(validation.Message), 400)
		}
	}

	return kubewarden.AcceptRequest()
}
