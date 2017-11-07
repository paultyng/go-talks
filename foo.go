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
			doTracedWork(ctx, "child1", 1, func(ctx context.Context) {
				doTracedWork(ctx, "child2", 2, func(ctx context.Context) {
					doTracedWork(ctx, "child3", 1, nil)
				})
			})
			return nil
		})

		g.Go(func() error {
			client := http.Client{
				Transport: &trace.Transport{}, // HL
			}
			req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d", barPort), nil)
			req = req.WithContext(ctx)

			_, err := client.Do(req)

			return err
		})

		err := g.Wait()
		if err != nil {
			log.Fatal(err)
		}
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", fooPort), tcli.HTTPHandler(handler))
}
