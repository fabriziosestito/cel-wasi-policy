package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/ext"
	k8s "github.com/kubewarden/go-wasi-policy-template/cel/library"
)

func validate(payload []byte) []byte {
	var validationRequest ValidationRequest

	err := json.Unmarshal(payload, &validationRequest)
	if err != nil {
		return marshalValidationResponseOrFail(
			RejectRequest(
				Message(fmt.Sprintf("Error deserializing validation request: %v", err)),
				Code(400)))
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

	ast, issues := env.Compile(validationRequest.Settings.Expression)
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

	if val == types.True {
		return marshalValidationResponseOrFail(
			AcceptRequest())
	}

	return marshalValidationResponseOrFail(
		RejectRequest(Message(validationRequest.Settings.Message), 400),
	)
}

func marshalValidationResponseOrFail(response ValidationResponse) []byte {
	responseBytes, err := json.Marshal(&response)
	if err != nil {
		log.Fatalf("cannot marshal validation response: %v", err)
	}
	return responseBytes
}
