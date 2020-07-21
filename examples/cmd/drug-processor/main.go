package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/google/uuid"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	eventsource "github.com/itsmurugappan/knative-eventing-sources/pkg/sources"
	"github.com/kelseyhightower/envconfig"
)

var (
	env            envConfig
	s3Client       *s3.S3
	anatomyDataMap map[string]*strings.Builder
	chunkSize      int
)

type envConfig struct {
	// s3 bucket
	Bucket string `envconfig:"S3_BUCKET"`

	Region string `envconfig:"S3_REGION"`

	Access string `envconfig:"S3_ACCESS_KEY"`

	Secret string `envconfig:"S3_SECRET_KEY"`

	Endpoint string `envconfig:"S3_URL"`

	// max file size to be uploaded
	FileSizeBuffer string `envconfig:"FILE_SIZE_BUFFER" default:"250000"`
}

func init() {
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}
	anatomyDataMap = make(map[string]*strings.Builder)
	chunkSize, _ = strconv.Atoi(env.FileSizeBuffer)

	newSession := session.New(&aws.Config{
		Region:      &env.Region,
		Endpoint:    &env.Endpoint,
		DisableSSL:  aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(env.Access, env.Secret, "")})

	s3Client = s3.New(newSession)
}

func main() {
	client, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Fatal("Failed to create client, ", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		writeData()
		cancel()
	}()

	log.Fatal(client.StartReceiver(ctx, process))
}

func writeData() {
	log.Println("writing data")
	for key := range anatomyDataMap {
		writeToStore(key)
	}
}

func process(event cloudevents.Event) {
	sourceData := eventsource.S3SourceData{}
	if err := event.DataAs(&sourceData); err != nil {
		fmt.Printf("got data error: %s\n", err.Error())
	}
	s := sourceData.Data

	for _, data := range strings.Split(s, "\n") {
		if data == "" {
			return
		}

		if data == "EOF" {
			log.Println("reached EOF")
			writeData()
			return
		}

		record := (strings.Split(data, ","))

		if len(record) < 4 {
			log.Printf("invalid data")
			continue
		}

		fileName := record[3]

		if fileName == "anatomical" {
			return
		}

		fileBuilder := anatomyDataMap[fileName]

		if fileBuilder == nil {
			fileBuilder = &strings.Builder{}
			anatomyDataMap[fileName] = fileBuilder
		}

		fileBuilder.WriteString(data)
		fileBuilder.WriteString("\n")

		if fileBuilder.Len() > chunkSize {
			writeToStore(fileName)
		}
	}
}

func writeToStore(fileName string) {
	fileBuilder := anatomyDataMap[fileName]
	if fileBuilder.Len() == 0 {
		//dont write
		return
	}
	s3FileName := fmt.Sprintf("%s-%s.csv", fileName, uuid.New())
	log.Printf("write to s3 %s with size %d", s3FileName, fileBuilder.Len())
	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(strings.NewReader(fileBuilder.String())),
		Bucket: aws.String(env.Bucket),
		Key:    aws.String(s3FileName),
	}

	_, err := s3Client.PutObject(input)
	if err != nil {
		log.Printf("error writing object to s3 %s", err.Error())
	}
	fileBuilder.Reset()
}
