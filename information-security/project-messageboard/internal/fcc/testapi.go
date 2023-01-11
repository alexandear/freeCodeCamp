package fcc

import (
	"encoding/json"
	"net/http"
	"strings"

	chi "github.com/go-chi/chi/v5"
)

func FCC() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r.Method == http.MethodGet && r.URL.Path == "/_api/app-info" {
					var header struct {
						Headers map[string]string `json:"headers"`
					}
					header.Headers = make(map[string]string)
					for name, values := range w.Header() {
						if (len(values)) > 0 {
							header.Headers[strings.ToLower(name)] = values[0]
						}
					}

					b, _ := json.Marshal(header)
					w.Header().Set("Content-Type", "application/json")
					w.Write(b)
				}
			}()

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func RegistersHandlers(r *chi.Mux) {
	r.Get("/_api/app-info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Get("/_api/get-tests", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
