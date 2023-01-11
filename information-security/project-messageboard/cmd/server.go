package cmd

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"time"

	env "github.com/caarlos0/env/v6"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Config struct {
	Port        int    `env:"PORT"`
	Environment string `env:"ENVIRONMENT"`
}

func ExecServer(embeddedFiles embed.FS) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Parse env: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.AllowAll().Handler)

	publicFS, err := fs.Sub(embeddedFiles, "public")
	if err != nil {
		log.Fatalf("sub public: %v", err)
	}
	r.Handle("/public/*", http.StripPrefix("/public", http.FileServer(http.FS(publicFS))))
	viewsFS, err := fs.Sub(embeddedFiles, "views")
	if err != nil {
		log.Fatalf("sub views: %v", err)
	}
	r.Handle("/*", http.FileServer(http.FS(viewsFS)))

	host := ""
	if cfg.Environment == "local" {
		host = "localhost"
	}
	serverAddr := host + ":" + strconv.Itoa(cfg.Port)
	s := http.Server{
		Addr:         serverAddr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Server is running on http://%s", serverAddr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("listen and serve: %v", err)
	}
}
