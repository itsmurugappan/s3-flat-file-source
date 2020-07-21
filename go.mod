module github.com/itsmurugappan/s3-flat-file-source

go 1.13

require (
	github.com/aws/aws-sdk-go v1.31.4
	github.com/cloudevents/sdk-go/v2 v2.0.0
	github.com/google/go-cmp v0.4.0
	github.com/google/uuid v1.1.1
	github.com/itsmurugappan/knative-eventing-sources v0.0.0-20200612162000-be1a29d5b241
	github.com/kelseyhightower/envconfig v1.4.0
	gomodules.xyz/jsonpatch/v2 v2.1.0 // indirect
	k8s.io/api v0.18.2 // indirect
	k8s.io/utils v0.0.0-20200414100711-2df71ebbae66 // indirect
)

replace k8s.io/client-go => k8s.io/client-go v0.17.4
