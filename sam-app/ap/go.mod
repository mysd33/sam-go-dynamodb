require (
	example.com/apbase v0.0.0-00010101000000-000000000000
	github.com/aws/aws-lambda-go v1.28.0
	github.com/aws/aws-sdk-go v1.42.15
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/pkg/errors v0.9.1
)

module ap

go 1.16

replace example.com/apbase => ../apbase
