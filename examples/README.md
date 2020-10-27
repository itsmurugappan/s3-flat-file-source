# Using S3 File Source

Examples on how to put the s3 file source to use. 
Both the example will process a drug list file.
File used is --> [here](./files/drug_uses.csv)

## Pre Req

Pre req's are common for the both the examples

1. Create the secrets for S3_ACCESS_KEY, S3_SECRET_KEY, S3_URL and S3_REGION 
```
kubectl create secret generic s3-access  \
--from-literal S3_REGION=us-west-1 \
--from-literal S3_URL=https://s3-us-west-1.amazonaws.com \
--from-literal S3_ACCESS_KEY= \
--from-literal S3_SECRET_KEY= -n demo
```
2. [Create rolebinding to give access to default service account to create a job](https://github.com/itsmurugappan/job-trigger#prereq)
    
  This step is to give access to the service account running ksvc to create a kubernetes job

3. [knative eventing core](https://knative.dev/docs/install/any-kubernetes-cluster/#installing-the-eventing-component)

4. Deploy `s3-file-source`

`s3-file-source` is a kubernetes-job, created by [job-trigger ksvc](https://github.com/itsmurugappan/job-trigger)

Apply the below yaml in your namespace.

```
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: s3-file-source-svc
  namespace: demo
spec:
  template:
    spec:
      containers:
      - env:
        - name: JOB_SPEC
          value: "{\"Image\": \"murugappans/s3-file-source:v1\",\"Name\": \"s3-file-source-job\", \"EnvFromSecretorCM\": [{\"Name\": \"s3-access\",\"Type\": \"Secret\"}]}"
        image: murugappans/job-trigger:v1
```

## Example 1 - Process a drug list and create separate files based on anatomy

This example will create separate files based on anatomy and upload it to s3 bucket.

```
kubectl apply -f config/drug-process.yaml -n demo
```

Above will deploy ksvc for processing drug file and sinkbinding for binding the source with this ksvc. 

To invoke, curl like below

```
curl "http://s3-file-source-job-trigger-demo.muru.fun/?labels=app=drug-file-source&S3_BUCKET=druglist&S3_FILE_NAME=drug_uses.csv"
```

## Example 2 - View events on cloud events viewer.

In this example we will just view the cloud events sent by the file source in a ui. 

```
kubectl apply -f config/view-ce.yaml -n demo
```

Above will deploy ksvc for view cloud events and sinkbinding for binding the source with this ksvc

Open the web page for {ksvc-url}/index.html

To invoke, curl like below

```
curl "http://s3-file-source-job-trigger-demo.muru.fun/?labels=app=view-drug-file-source&S3_BUCKET=druglist&S3_FILE_NAME=drug_uses.csv"
```
