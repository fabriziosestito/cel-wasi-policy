package validate

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/google/cel-go/common/types"
	"github.com/kubewarden/cel-policy/internal/cel"
	"github.com/kubewarden/cel-policy/internal/settings"
	kubewarden "github.com/kubewarden/policy-sdk-go"
	"github.com/kubewarden/policy-sdk-go/pkg/capabilities"
	"github.com/kubewarden/policy-sdk-go/pkg/capabilities/kubernetes"
	kubewardenProtocol "github.com/kubewarden/policy-sdk-go/protocol"
)

var host = capabilities.NewHost()

func Validate(payload []byte) ([]byte, error) {
	validationRequest := kubewardenProtocol.ValidationRequest{}

	err := json.Unmarshal(payload, &validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(fmt.Sprintf("Error deserializing validation request: %v", err)),
			kubewarden.Code(400))
	}

	settings, err := settings.NewSettingsFromValidationReq(&validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(fmt.Sprintf("Error serializing RawMessage: %v", err)),
			kubewarden.Code(400))
	}

	compiler, err := cel.NewCompiler()
	if err != nil {
		return nil, fmt.Errorf("failed to create CEL compiler: %w", err)
	}

	object := map[string]interface{}{}
	err = json.Unmarshal(validationRequest.Request.Object, &object)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal request object %w", err)
	}

	oldObject := map[string]interface{}{}
	if validationRequest.Request.OldObject != nil {
		err = json.Unmarshal(validationRequest.Request.OldObject, &oldObject)
		if err != nil {
			return nil, fmt.Errorf("cannot unmarshal request oldObject %w", err)
		}
	}

	// TODO: change this to corev1.Namespace once nested structs are supported
	// See: https://github.com/google/cel-go/pull/892
	var namespaceObject map[string]interface{}
	if validationRequest.Request.Namespace != "" {
		resourceRequest := kubernetes.GetResourceRequest{
			APIVersion: "v1",
			Kind:       "Namespace",
			Name:       validationRequest.Request.Namespace,
		}

		responseBytes, err := kubernetes.GetResource(&host, resourceRequest)
		if err != nil {
			return nil, fmt.Errorf("cannot get namespace data: %w", err)
		}
		if err := json.Unmarshal(responseBytes, &namespaceObject); err != nil {
			return nil, fmt.Errorf("cannot parse namespace data: %w", err)
		}
	}

	vars := map[string]interface{}{
		"object":          object,
		"oldObject":       oldObject,
		"request":         validationRequest.Request,
		"namespaceObject": namespaceObject,
	}

	if err := evalVariables(compiler, vars, settings.Variables); err != nil {
		return nil, fmt.Errorf("failed to evaluate variables: %w", err)
	}

	return evalValidations(compiler, vars, settings.Validations)
}

func evalVariables(compiler *cel.Compiler, vars map[string]interface{}, variables []settings.Variable) error {
	for _, variable := range variables {
		ast, err := compiler.CompileCELExpression(variable.Expression)
		if err != nil {
			return err
		}

		val, err := compiler.EvalCELExpression(vars, ast)
		if err != nil {
			return err
		}

		if err := compiler.AddVariable(variable.Name, ast.OutputType()); err != nil {
			return err
		}
		vars[fmt.Sprintf("variables.%s", variable.Name)] = val
	}
	return nil
}

func evalValidations(compiler *cel.Compiler, vars map[string]interface{}, validations []settings.Validation) ([]byte, error) {
	for _, validation := range validations {
		ast, err := compiler.CompileCELExpression(validation.Expression)
		if err != nil {
			return nil, fmt.Errorf("failed to compile expression: %w", err)
		}

		val, err := compiler.EvalCELExpression(vars, ast)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate expression: %w", err)
		}

		if val == types.False {
			reason := reasonToStatusCode(validation.Reason)

			if validation.MessageExpression != "" {
				message, err := evalMessageExpression(compiler, vars, validation.MessageExpression)
				if err != nil {
					return nil, fmt.Errorf("failed to evaluate message expression: %w", err)
				}
				return kubewarden.RejectRequest(kubewarden.Message(message), reason)
			}
			if validation.Message != "" {
				return kubewarden.RejectRequest(kubewarden.Message(validation.Message), reason)
			}

			return kubewarden.RejectRequest(kubewarden.Message(fmt.Sprintf("failed expression: %s", strings.TrimSpace(validation.Expression))), reason)
		}
	}

	return kubewarden.AcceptRequest()
}

func evalMessageExpression(compiler *cel.Compiler, vars map[string]interface{}, messageExpression string) (string, error) {
	ast, err := compiler.CompileCELExpression(messageExpression)
	if err != nil {
		return "", err
	}

	val, err := compiler.EvalCELExpression(vars, ast)
	if err != nil {
		return "", err
	}

	message, ok := val.Value().(string)
	if !ok {
		return "", fmt.Errorf("message expression must evaluate to string")
	}

	return message, nil
}

func reasonToStatusCode(reason string) kubewarden.Code {
	var statusCode kubewarden.Code
	switch reason {
	case settings.StatusReasonInvalid:
		statusCode = kubewarden.Code(400)
	case settings.StatusReasonForbidden:
		statusCode = 403
	case settings.StatusReasonRequestEntityTooLarge:
		statusCode = 413
	case settings.StatusReasonUnauthorized:
		statusCode = 401
	default:
		// This should never happen since we validate the settings when loading the policy
		log.Fatalf("unknown reason: %v", reason)
	}

	return statusCode
}
