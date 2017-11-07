package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"cloud.google.com/go/trace"
)

func serveFoo(tcli *trace.Client) error {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("handling foo request")
		ctx := r.Context()

		if span := trace.FromContext(ctx); span != nil {
			span.SetLabel("user_agent", r.UserAgent())
		}

		g, ctx := errgroup.WithContext(ctx)

		// a little prep work
		doWork(1)

		g.Go(func() error {
			// some long running task
			doTracedWork(ctx, 5)
			return nil
		})

		g.Go(func() error {
			barClient := http.Client{
				Transport: &trace.Transport{},
			}
			req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d", barPort), nil)
			req = req.WithContext(ctx)

			_, err := barClient.Do(req)

			return err
		})

		err := g.Wait()
		if err != nil {
			log.Fatal(err)
		}
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", fooPort), tcli.HTTPHandler(handler))
}

func doTracedWork(ctx context.Context, seconds int) {
	if parentSpan := trace.FromContext(ctx); parentSpan != nil {
		childSpan := parentSpan.NewChild("doSomeLongRunningThing")
		defer childSpan.Finish()
	}

	doWork(seconds)
}

func doWork(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
