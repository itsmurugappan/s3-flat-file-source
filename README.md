# S3 File Source 

Emit lines of flat files stored in s3 based object storage as cloud events

[![Go Report Card](https://goreportcard.com/badge/github.com/itsmurugappan/s3-flat-file-source)](https://goreportcard.com/report/github.com/itsmurugappan/s3-flat-file-source)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-59%25-brightgreen.svg?longCache=true&style=flat)</a>


![](./images/fileprocessor.jpg)

### Configuration

Following are the configuration options (Environment Variables) for the file source 

1. S3_URL
2. S3_BUCKET
3. S3_REGION
4. S3_ACCESS_KEY
5. S3_SECRET_KEY
6. S3_FILE_NAME
7. CHUNK_SIZE - Bytes to be downloaded each time from s3 (Default - 50mb)
8. SINK_DUMP_COUNT - Number of lines to be sent to the sink in one transaction (Default - 100)
9. SINK_RETRY_COUNT - Number of times to retry on failure (Default - 3)
10. SINK_RETRY_INTERVAL - Seconds to wait before next retry (Default - 1)

### Setting Up the File Source

This app will run to completion, so I have leveraged the [Knative-Job Trigger](https://github.com/itsmurugappan/job-trigger) 
to get a serverless webhook to deploy the job dyanmically.

#### Pre Req

1. Create the secrets for S3_ACCESS_KEY, S3_SECRET_KEY, S3_URL and S3_REGION 
2. [Create rolebinding to give access to default service account to create a job](https://github.com/itsmurugappan/job-trigger#prereq)

#### Creating the Knative Service

```
kubectl apply -f - <<EOF
apiVersion: serving.knative.dev/v1alpha1
kind: Service
metadata:
  name: s3-file-source-job-trigger
spec:
  template:
    spec:
      containers:
      - env:
        - name: spec
          value: "{\"Image\": \"docker.io/murugappans/s3-file-source-8baccff8d26c20f105867f48bc2124cf:v1\",\"Name\": \"s3-file-source-job\", \"EnvFromSecretorCM\": [{\"Name\": \"region\",\"Type\": \"Secret\"},{\"Name\": \"secret\",\"Type\": \"Secret\"},{\"Name\": \"access\",\"Type\": \"Secret\"},{\"Name\": \"endpoint\",\"Type\": \"Secret\"}]}"
        image: ko://github.com/itsmurugappan/cmd/job-trigger
EOF
```

Now the ksvc should up and running. To trigger the job, ksvc needs to be invoked.
Before that, the sink binding and the event processor needs to be set up.

### Processing the Cloud Events

The example folder has an example of file processing.
1. Process a generic drug csv and organises the files based on anatomy

Running the drug file example

```
kubectl apply -f examples/cmd/drug-processor/ksvc.yaml
```

### Generate Events

Now that all components are ready, invoke the file source service to generate events

```
curl -X POST 'http://s3-file-source-job-trigger-test-alpha.faas.com/?labels=app=drug-file-source&S3_BUCKET=drug-raw-files&S3_FILE_NAME=drug_uses.csv'
```
