package cel

const (
	// PerCallLimit specify the actual cost limit per CEL validation call
	// current PerCallLimit gives roughly 0.1 second for each expression validation call
	PerCallLimit = 1000000

	// RuntimeCELCostBudget is the overall cost budget for runtime CEL validation cost per ValidatingAdmissionPolicyBinding or CustomResource
	// current RuntimeCELCostBudget gives roughly 1 seconds for the validation
	RuntimeCELCostBudget = 10000000

	// RuntimeCELCostBudgetMatchConditions is the overall cost budget for runtime CEL validation cost on matchConditions per object with matchConditions
	// this is per webhook for validatingwebhookconfigurations and mutatingwebhookconfigurations or per ValidatingAdmissionPolicyBinding
	// current RuntimeCELCostBudgetMatchConditions gives roughly 1/4 seconds for the validation
	RuntimeCELCostBudgetMatchConditions = 2500000

	// CheckFrequency configures the number of iterations within a comprehension to evaluate
	// before checking whether the function evaluation has been interrupted
	CheckFrequency = 100

	// MaxRequestSizeBytes is the maximum size of a request to the API server
	// TODO(DangerOnTheRanger): wire in MaxRequestBodyBytes from apiserver/pkg/server/options/server_run_options.go to make this configurable
	// Note that even if server_run_options.go becomes configurable in the future, this cost constant should be fixed and it should be the max allowed request size for the server
	MaxRequestSizeBytes = int64(3 * 1024 * 1024)

	// MaxEvaluatedMessageExpressionSizeBytes represents the largest-allowable string generated
	// by a messageExpression field
	MaxEvaluatedMessageExpressionSizeBytes = 5 * 1024
)

const (
	// DefaultMaxRequestSizeBytes is the size of the largest request that will be accepted
	DefaultMaxRequestSizeBytes = MaxRequestSizeBytes

	// MaxDurationSizeJSON
	// OpenAPI duration strings follow RFC 3339, section 5.6 - see the comment on maxDatetimeSizeJSON
	MaxDurationSizeJSON = 32
	// MaxDatetimeSizeJSON
	// OpenAPI datetime strings follow RFC 3339, section 5.6, and the longest possible
	// such string is 9999-12-31T23:59:59.999999999Z, which has length 30 - we add 2
	// to allow for quotation marks
	MaxDatetimeSizeJSON = 32
	// MinDurationSizeJSON
	// Golang allows a string of 0 to be parsed as a duration, so that plus 2 to account for
	// quotation marks makes 3
	MinDurationSizeJSON = 3
	// JSONDateSize is the size of a date serialized as part of a JSON object
	// RFC 3339 dates require YYYY-MM-DD, and then we add 2 to allow for quotation marks
	JSONDateSize = 12
	// MinDatetimeSizeJSON is the minimal length of a datetime formatted as RFC 3339
	// RFC 3339 datetimes require a full date (YYYY-MM-DD) and full time (HH:MM:SS), and we add 3 for
	// quotation marks like always in addition to the capital T that separates the date and time
	MinDatetimeSizeJSON = 21
	// MinStringSize is the size of literal ""
	MinStringSize = 2
	// MinBoolSize is the length of literal true
	MinBoolSize = 4
	// MinNumberSize is the length of literal 0
	MinNumberSize = 1
)
