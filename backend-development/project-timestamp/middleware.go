package main

import (
	"log"
	"net/http"
)

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.RemoteAddr, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}
