package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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

	serverAddr := host + ":" + port
	if err := r.Run(serverAddr); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
