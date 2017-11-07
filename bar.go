package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/trace"
	"github.com/ExpansiveWorlds/instrumentedsql"
	"github.com/ExpansiveWorlds/instrumentedsql/google"
	"github.com/lib/pq"
)

func init() {
	sql.Register("postgres-trace", instrumentedsql.WrapDriver(&pq.Driver{}, instrumentedsql.WithTracer(google.NewTracer())))
}

func serveBar(tcli *trace.Client) error {
	db, err := sql.Open("postgres-trace", "host=localhost dbname=postgres user=postgres sslmode=disable")
	if err != nil {
		return err
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("handling bar request")

		ctx := r.Context()

		//prep work
		doWork(1)

		_, err := db.ExecContext(ctx, "select pg_sleep(2)")
		if err != nil {
			log.Fatal(err)
		}

		//process query results
		doWork(4)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", barPort), tcli.HTTPHandler(handler))
}
