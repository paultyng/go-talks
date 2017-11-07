package main

import (
	"context"
	"flag"
	"log"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/trace"
)

var mode string
var projectID string

const (
	fooPort = 3000
	barPort = 3001
)

func init() {
	flag.StringVar(&mode, "mode", "foo", "runtime mode: `foo` or `bar`")
	flag.StringVar(&projectID, "project", "", "GCP project ID")
}

func main() {
	flag.Parse()

	if projectID == "" {
		if metadata.OnGCE() {
			if pid, err := metadata.ProjectID(); err == nil {
				projectID = pid
			}
		}
	}

	log.Printf("project %s", projectID)

	ctx := context.Background()
	tcli, err := trace.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}
	pol, err := trace.NewLimitedSampler(1, 100)
	if err != nil {
		log.Fatal(err)
	}
	tcli.SetSamplingPolicy(pol)

	log.Printf("starting '%s'", mode)

	switch mode {
	case "foo":
		log.Fatal(serveFoo(tcli))
	case "bar":
		log.Fatal(serveBar(tcli))
	default:
		log.Fatalf("mode '%s' not expected", mode)
	}
}
