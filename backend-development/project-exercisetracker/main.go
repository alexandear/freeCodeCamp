package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongodbName = "exercise_tracker"
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

	mongoURI := os.Getenv("MONGODB_URI")

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPIOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	e := echo.New()
	mongoDB := mongoClient.Database(mongodbName)
	userServ := newUserService(mongoDB)
	h := newHandler(e, userServ)
	serverAddr := host + ":" + port
	s := http.Server{
		Addr:         serverAddr,
		Handler:      h,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Printf("Server is running on http://%s", serverAddr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
