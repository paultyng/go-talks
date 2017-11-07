package main

import (
	"log"
	"net/http"
	"net/http/httputil"
)

func newRequestLoggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := httputil.DumpRequest(r, false)
		log.Println(string(data))
		next.ServeHTTP(w, r)
	})
}
