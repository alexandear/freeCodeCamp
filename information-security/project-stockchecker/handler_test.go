package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
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
	var stockPrice StockPriceResp
	if err := json.NewDecoder(bytes.NewReader(resBytes)).Decode(&stockPrice); err != nil {
		t.Fatal(err)
	}

	if stockPrice.StockData.Price == 0.0 {
		t.Fatalf("price must be non-zero")
	}

	expected := fmt.Sprintf(`{"stockData":{"stock":"GRMN","price":%g,"likes":0}}
`, stockPrice.StockData.Price)
	if expected != actual {
		t.Fatalf("expected %+v, got %+v", expected, actual)
	}
}

func TestHandler_StockPrice_GetOneStockAndLike(t *testing.T) {
	h := NewHandler(echo.New(), NewStockService(db))

	s := httptest.NewServer(h)
	defer s.Close()

	res, err := client.Get(s.URL + "/api/stock-prices?stock=AAPL&like=true")
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
	var stockPrice StockPriceResp
	if err := json.NewDecoder(bytes.NewReader(resBytes)).Decode(&stockPrice); err != nil {
		t.Fatal(err)
	}

	expected := fmt.Sprintf(`{"stockData":{"stock":"AAPL","price":%g,"likes":1}}
`, stockPrice.StockData.Price)
	if expected != actual {
		t.Fatalf("expected %+v, got %+v", expected, actual)
	}
}

func TestHandler_StockPrice_GetOneStockAndLikeFewTimes(t *testing.T) {
	h := NewHandler(echo.New(), NewStockService(db))

	s := httptest.NewServer(h)
	defer s.Close()

	if _, err := client.Get(s.URL + "/api/stock-prices?stock=MSFT&like=true"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Get(s.URL + "/api/stock-prices?stock=MSFT&like=true"); err != nil {
		t.Fatal(err)
	}
	res, err := client.Get(s.URL + "/api/stock-prices?stock=MSFT&like=true")
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
	var stockPrice StockPriceResp
	if err := json.NewDecoder(bytes.NewReader(resBytes)).Decode(&stockPrice); err != nil {
		t.Fatal(err)
	}

	expected := fmt.Sprintf(`{"stockData":{"stock":"MSFT","price":%g,"likes":1}}
`, stockPrice.StockData.Price)
	if expected != actual {
		t.Fatalf("expected %+v, got %+v", expected, actual)
	}
}

func TestHandler_StockPrice_GetTwoStocks(t *testing.T) {
	h := NewHandler(echo.New(), NewStockService(db))

	s := httptest.NewServer(h)
	defer s.Close()

	res, err := client.Get(s.URL + "/api/stock-prices?stock=GOOG&stock=MSFT")
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
	var stockPrices StockPricesResp
	if err := json.NewDecoder(bytes.NewReader(resBytes)).Decode(&stockPrices); err != nil {
		t.Fatal(err)
	}

	expected := fmt.Sprintf(`{"stockData":[{"stock":"GOOG","price":%g,"rel_likes":0},{"stock":"MSFT","price":%g,"rel_likes":0}]}
`, stockPrices.StockData[0].Price, stockPrices.StockData[1].Price)
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

	res := mongoClient.Database("test_stock_checker")

	_, _ = res.Collection("stocks").DeleteMany(context.Background(), bson.D{{}})
	_, _ = res.Collection("stock_per_ips").DeleteMany(context.Background(), bson.D{{}})

	return res
}
