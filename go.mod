module github.com/itsmurugappan/s3-flat-file-source

go 1.15

require (
	github.com/aws/aws-sdk-go v1.31.12
	github.com/cloudevents/sdk-go/v2 v2.0.0
	github.com/google/go-cmp v0.5.2
	github.com/google/uuid v1.1.1
	github.com/itsmurugappan/knative-eventing-sources v0.0.0-20201026190751-81ef2f2fe1eb
	github.com/kelseyhightower/envconfig v1.4.0
	gomodules.xyz/jsonpatch/v2 v2.1.0 // indirect
)

replace k8s.io/client-go => k8s.io/client-go v0.18.8
