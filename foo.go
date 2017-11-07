package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"

	"cloud.google.com/go/trace"
)

func serveFoo(tcli *trace.Client) error {
	// create handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("handling foo request")
		ctx := r.Context()

		if span := trace.FromContext(ctx); span != nil {
			// set a sample label on the span if present
			span.SetLabel("user_agent", r.UserAgent())
		}

		// a little prep work
		doWork(1)

		// start an errgroup for parallel work
		g, ctx := errgroup.WithContext(ctx)

		g.Go(func() error {
			// some long running nested tasks
			doTracedWork(ctx, "child1", 1, func(ctx context.Context) {
				doTracedWork(ctx, "child2", 2, func(ctx context.Context) {
					doTracedWork(ctx, "child3", 1, nil)
				})
			})
			return nil
		})

		g.Go(func() error {
			client := http.Client{
				Transport: &trace.Transport{},
			}
			req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d", barPort), nil)
			req = req.WithContext(ctx)

			// execute a remote request against the "bar" service propagating tracing
			_, err := client.Do(req)

			return err
		})

		// wait for the errgroup
		err := g.Wait()
		if err != nil {
			log.Fatal(err)
		}
	})

	// wrap handler with tracing
	return http.ListenAndServe(fmt.Sprintf(":%d", fooPort), tcli.HTTPHandler(handler))
}
