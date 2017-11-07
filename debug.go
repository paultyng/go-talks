package main

import (
	"net/http"
	"os"
)

func newRequestLoggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Write(os.Stdout)
	})
}
