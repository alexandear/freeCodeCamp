package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"stockchecker/fcc"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	port := "0"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	host := ""
	if os.Getenv("ENVIRONMENT") == "local" {
		host = "localhost"
	}

	e := echo.New()
	e.Use(middleware.CORS()) // for testing purposes only
	e.Use(fcc.FCC())
	e.Static("/", "public")
	e.File("/", "views/index.html")
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		ContentSecurityPolicy: "default-src 'self'; script-src 'self'; style-src 'self'",
	}))

	serverAddr := host + ":" + port
	s := http.Server{
		Addr:         serverAddr,
		Handler:      e,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Printf("Server is running on http://%s", serverAddr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
