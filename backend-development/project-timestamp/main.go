package main

import (
	"embed"
	"log"
	"net/http"
	"strings"

	"github.com/rs/cors"
)

//go:embed views
//go:embed public
var staticFiles embed.FS

const serverAddr = ":0"

func rootPath(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			r.URL.Path = "/views/"
		} else {
			b := strings.Split(r.URL.Path, "/")[0]
			if b != "public" {
				r.URL.Path = "/public" + r.URL.Path
			}
		}
		h.ServeHTTP(w, r)
	})
}

func main() {
	staticFS := http.FS(staticFiles)
	fs := rootPath(http.FileServer(staticFS))

	api := apiHandler{
		clock: &realClock{},
	}

	mux := http.NewServeMux()
	mux.Handle("/", fs)
	mux.Handle("/api", api)
	mux.Handle("/api/", api)
	handler := cors.AllowAll().Handler(loggerMiddleware(mux))

	log.Printf("Server is running on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, handler))
}
