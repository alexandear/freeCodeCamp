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

	"messageboard/api"
	httpserv "messageboard/http"
	"messageboard/internal/fcc"
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
	r.Use(cors.AllowAll().Handler, fcc.FCC()) // for testing purposes
	r.Use(middleware.SetHeader("X-DNS-Prefetch-Control", "off"))
	r.Use(middleware.SetHeader("X-Frame-Options", "SAMEORIGIN"))
	r.Use(middleware.SetHeader("Referrer-Policy", "same-origin"))
	fcc.RegistersHandlers(r)

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

	serv := httpserv.NewServer()
	api.HandlerFromMux(serv, r)

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
