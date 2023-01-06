package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/rs/cors"
)

//go:embed views
//go:embed public
var staticFiles embed.FS

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
	port := "0"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	host := ""
	if os.Getenv("ENVIRONMENT") == "local" {
		host = "localhost"
	}

	staticFS := http.FS(staticFiles)
	fs := rootPath(http.FileServer(staticFS))

	mux := http.NewServeMux()
	mux.Handle("/", fs)
	handler := cors.AllowAll().Handler(mux)

	serverAddr := host + ":" + port
	log.Printf("Server is running on http://%s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, handler))
}
