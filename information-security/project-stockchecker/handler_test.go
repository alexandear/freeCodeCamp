package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client = http.Client{Timeout: time.Second}
	db     = newTestMongoDatabase()
)

func TestHandler_StockPrice_GetOneStock(t *testing.T) {
	h := NewHandler(echo.New(), NewStockService(db))

	s := httptest.NewServer(h)
	defer s.Close()

	res, err := client.Get(s.URL + "/api/stock-prices?stock=GRMN")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if http.StatusOK != res.StatusCode {
		t.Fatalf("expected '200 OK' status, got '%s'", res.Status)
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	actual := string(resBytes)

	expected := `{"stockData":{"stock":"GRMN","price":98.33,"likes":0}}
`
	if expected != actual {
		t.Fatalf("expected %+v, got %+v", expected, actual)
	}
}

func newTestMongoDatabase() *mongo.Database {
	mongoURI := os.Getenv("MONGODB_URI")

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPIOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	return mongoClient.Database("test_stock_checker")
}
