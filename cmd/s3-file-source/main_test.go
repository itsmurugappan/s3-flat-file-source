package main

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	eventsource "github.com/itsmurugappan/knative-eventing-sources/pkg/sources"
	eventsourcetest "github.com/itsmurugappan/knative-eventing-sources/pkg/test/sources"
)

// test to see if env variables are getting parsed.
func TestConstructS3Source(t *testing.T) {
	opt := cmpopts.IgnoreUnexported(eventsource.S3Source{})

	os.Setenv("K_SINK", "http://s3processor.com")
	os.Setenv("SINK_DUMP_COUNT", "10")
	os.Setenv("DOWNLOAD_CHUNK_SIZE", "100000000")
	os.Setenv("S3_BUCKET", "test-bucket")
	os.Setenv("S3_FILE_NAME", "test.csv")
	os.Setenv("S3_REGION", "us-west")
	os.Setenv("S3_ACCESS_KEY", "ACCESS-KEY")
	os.Setenv("S3_SECRET_KEY", "ACCESS-SECRET")
	os.Setenv("S3_URL", "https://s3url.com")
	os.Setenv("SINK_RETRY_COUNT", "5")
	os.Setenv("SINK_RETRY_INTERVAL", "1")

	want := &eventsource.S3Source{
		Sink: eventsource.SinkConfig{
			URL:           "http://s3processor.com",
			DumpCount:     "10",
			Retries:       "5",
			RetryInterval: "1",
		},
		Connection: eventsource.ObjectStoreConfig{
			Bucket:    "test-bucket",
			Key:       "test.csv",
			Region:    "us-west",
			Access:    "ACCESS-KEY",
			Secret:    "ACCESS-SECRET",
			Endpoint:  "https://s3url.com",
			ChunkSize: "100000000",
		},
	}
	if diff := cmp.Diff(want, constructS3Source(), opt); diff != "" {
		t.Errorf("ConstructS3Source() mismatch (-want +got):\n%s", diff)
	}
}

func TestStartSource(t *testing.T) {
	input := &eventsourcetest.S3SourceFake{}
	gotI := startEventSource(input)
	got, _ := gotI.(string)
	if diff := cmp.Diff("Events Generated", got); diff != "" {
		t.Errorf("StartEventSource() mismatch (-want +got):\n%s", diff)
	}
}
