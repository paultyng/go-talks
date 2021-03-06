package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"cloud.google.com/go/trace"
	"github.com/ExpansiveWorlds/instrumentedsql"
	"github.com/ExpansiveWorlds/instrumentedsql/google"
	"github.com/lib/pq"
)

func init() {
	// register instrumented sql driver that wraps postgres
	sql.Register("postgres-trace", instrumentedsql.WrapDriver(&pq.Driver{}, instrumentedsql.WithTracer(google.NewTracer())))
}

func serveBar(tcli *trace.Client) error {
	// connect to db
	db, err := sql.Open("postgres-trace", "host=localhost dbname=postgres user=postgres sslmode=disable")
	if err != nil {
		return err
	}

	// create handler that does some db work
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("handling bar request")

		ctx := r.Context()

		//prep work
		doWork(1)

		sleep := (1000.0 + float32(rand.Int31n(2000))) / 1000.0
		// ExecContext!!! very important for tracing
		_, err := db.ExecContext(ctx, "select pg_sleep($1)", sleep)
		if err != nil {
			log.Fatal(err)
		}

		// pretend to process query results
		doWork(2)
	})

	// wrap handler with request logger and tracing
	return http.ListenAndServe(fmt.Sprintf(":%d", barPort), newRequestLoggingHandler(tcli.HTTPHandler(handler)))
}
