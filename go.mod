module github.com/kubewarden/go-wasi-policy-template

go 1.21

require (
	github.com/google/cel-go v0.17.7
	github.com/kubewarden/policy-sdk-go v0.6.0
	google.golang.org/genproto/googleapis/api v0.0.0-20230803162519-f966b187b2e5
	google.golang.org/protobuf v1.31.0
)

replace github.com/go-openapi/strfmt => github.com/kubewarden/strfmt v0.1.3

require (
	github.com/antlr/antlr4/runtime/Go/antlr/v4 v4.0.0-20230305170008-8188dc5388df // indirect
	github.com/go-openapi/strfmt v0.21.3 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/kubewarden/k8s-objects v1.27.0-kw4 // indirect
	github.com/stoewer/go-strcase v1.2.0 // indirect
	github.com/wapc/wapc-guest-tinygo v0.3.3 // indirect
	golang.org/x/exp v0.0.0-20230515195305-f3d0a9c9a5cc // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
)
