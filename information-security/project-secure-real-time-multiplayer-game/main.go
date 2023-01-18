package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"

	"secure-real-time-multiplayer-game/internal/fcc"
)

func main() {
	port := "0"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	ginMode := gin.ReleaseMode
	host := ""
	if os.Getenv("ENVIRONMENT") == "local" {
		host = "localhost"
		ginMode = gin.DebugMode
	}

	gin.SetMode(ginMode)

	r := gin.Default()
	r.Static("/public", "./public")
	r.StaticFile("/", "./views/index.html")
	r.Use(cors.Default(), fcc.FCC()) // for testing

	r.Use(func(c *gin.Context) {
		c.Header("Surrogate-Control", "no-store")
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		c.Header("X-Powered-By", "PHP 7.4.3")
	}, secure.New(secure.Config{
		BrowserXssFilter:   true,
		ContentTypeNosniff: true,
	}))

	serverAddr := host + ":" + port
	if err := r.Run(serverAddr); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
