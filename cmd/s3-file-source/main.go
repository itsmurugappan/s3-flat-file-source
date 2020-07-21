package main

import (
	"flag"
	"log"
	"os"

	eventsource "github.com/itsmurugappan/knative-eventing-sources/pkg/sources"
	"github.com/kelseyhightower/envconfig"
)

var (
	sink string
)

func init() {
	flag.StringVar(&sink, "sink", "", "the host url")
}

func main() {
	flag.Parse()
	resultInterface := startEventSource(constructS3Source())
	result, _ := resultInterface.(eventsource.Result)
	log.Printf("Events Processing Stats \n Error Count: %d \n Success Count: %d", result.ErrorCount, result.SentCount)
}

func startEventSource(source eventsource.EventSource) interface{} {
	return eventsource.SourceEvents(source)
}

func constructS3Source() *eventsource.S3Source {
	var objectStoreMeta eventsource.ObjectStoreConfig
	var sinkConfig eventsource.SinkConfig

	if err := envconfig.Process("", &objectStoreMeta); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}
	if err := envconfig.Process("", &sinkConfig); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}

	if sink != "" {
		sinkConfig.URL = sink
	}

	return &eventsource.S3Source{Sink: sinkConfig, Connection: objectStoreMeta}
}
