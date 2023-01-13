package cmd

import (
	"context"
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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"messageboard/api"
	"messageboard/httpserv"
	"messageboard/internal/fcc"
	"messageboard/internal/gotest"
	"messageboard/msgboard"
)

type Config struct {
	Port        int    `env:"PORT"`
	Environment string `env:"ENVIRONMENT"`

	MongodbURI  string `env:"MONGODB_URI"`
	MongodbName string `env:"MONGODB_NAME" envDefault:"message_board"`

	RunTestsOnStart bool `env:"RUN_TESTS_ON_START" envDefault:"false"`
}

func ExecServer(embeddedFiles embed.FS) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Parse env: %v", err)
	}

	var tr *gotest.TestResults
	if cfg.RunTestsOnStart {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		res, err := gotest.Run(ctx, "test", nil, true)
		if err != nil {
			log.Fatalf("Run failed")
		}
		tr = res
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.AllowAll().Handler, fcc.FCC(tr)) // for testing purposes
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

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(cfg.MongodbURI).SetServerAPIOptions(serverAPIOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("mongo client connect: %v", err)
	}
	defer func() {
		dCtx, dCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer dCancel()

		if err := mongoClient.Disconnect(dCtx); err != nil {
			log.Fatal(err)
		}
	}()

	mongoDB := mongoClient.Database(cfg.MongodbName)

	msgServ := msgboard.NewService(mongoDB)
	serv := httpserv.NewServer(msgServ)
	strictHandler := api.NewStrictHandler(serv, nil)
	api.HandlerFromMux(strictHandler, r)

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
