package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"stockchecker/fcc"
	"stockchecker/gotest"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongodbName = "stock_checker"
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

	var tr *gotest.TestResults
	if os.Getenv("RUN_TESTS_ON_START") == "true" {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		res, err := gotest.Run(ctx, ".", nil, true)
		if err != nil {
			log.Fatalf("Run failed")
		}
		tr = res
	}

	mongoURI := os.Getenv("MONGODB_URI")

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPIOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		dctx, dcancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer dcancel()

		if err := mongoClient.Disconnect(dctx); err != nil {
			log.Fatal(err)
		}
	}()

	mongoDB := mongoClient.Database(mongodbName)
	stock := NewStockService(mongoDB)
	e := echo.New()
	h := NewHandler(e, stock)

	e.Use(middleware.CORS()) // for testing purposes only
	e.Use(fcc.FCC(tr))
	e.Static("/", "public")
	e.File("/", "views/index.html")
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		ContentSecurityPolicy: "default-src 'self'; script-src 'self'; style-src 'self'",
	}))

	serverAddr := host + ":" + port
	s := http.Server{
		Addr:         serverAddr,
		Handler:      h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Printf("Server is running on http://%s", serverAddr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
