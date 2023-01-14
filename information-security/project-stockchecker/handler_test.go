package main

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestHandler_StockPrice_GetOneStock(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	var actual StockPriceResp
	err := requests.URL(s.URL + "/api/stock-prices?stock=GRMN").ToJSON(&actual).Fetch(context.Background())

	require.NoError(t, err)
	assert.NotZero(t, actual.StockData.Price)
	expected := StockPriceResp{
		StockData: StockDataResp{
			Stock: "GRMN",
			Price: actual.StockData.Price,
			Likes: 0,
		},
	}
	assert.Equal(t, expected, actual)
}

func TestHandler_StockPrice_GetOneStockAndLike(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	var actual StockPriceResp
	err := requests.URL(s.URL + "/api/stock-prices?stock=AAPL&like=true").ToJSON(&actual).Fetch(context.Background())

	require.NoError(t, err)
	assert.NotZero(t, actual.StockData.Price)
	expected := StockPriceResp{
		StockData: StockDataResp{
			Stock: "AAPL",
			Price: actual.StockData.Price,
			Likes: 1,
		},
	}
	assert.Equal(t, expected, actual)
}

func TestHandler_StockPrice_GetOneStockAndLikeFewTimes(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	err := requests.URL(s.URL + "/api/stock-prices?stock=MSFT&like=true").Fetch(context.Background())
	require.NoError(t, err)
	err = requests.URL(s.URL + "/api/stock-prices?stock=MSFT&like=true").Fetch(context.Background())
	require.NoError(t, err)

	var actual StockPriceResp
	err = requests.URL(s.URL + "/api/stock-prices?stock=MSFT&like=true").ToJSON(&actual).Fetch(context.Background())

	require.NoError(t, err)
	assert.NotZero(t, actual.StockData.Price)
	expected := StockPriceResp{
		StockData: StockDataResp{
			Stock: "MSFT",
			Price: actual.StockData.Price,
			Likes: 1,
		},
	}
	assert.Equal(t, expected, actual)
}

func TestHandler_StockPrice_GetTwoStocks(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	var actual StockPricesResp
	err := requests.URL(s.URL + "/api/stock-prices?stock=META&stock=INTC").ToJSON(&actual).Fetch(context.Background())

	require.NoError(t, err)
	require.Len(t, actual.StockData, 2)
	assert.NotZero(t, actual.StockData[0].Price)
	assert.NotZero(t, actual.StockData[1].Price)

	expected := StockPricesResp{
		StockData: []StockDatasResp{
			{
				Stock:    "META",
				Price:    actual.StockData[0].Price,
				RelLikes: 0,
			},
			{
				Stock:    "INTC",
				Price:    actual.StockData[1].Price,
				RelLikes: 0,
			},
		},
	}
	assert.Equal(t, expected, actual)
}

func TestHandler_StockPrice_TwoStocksWithLikes(t *testing.T) {
	stockServ := NewStockService(newTestMongoDatabase())
	h := NewHandler(echo.New(), stockServ)
	s := httptest.NewServer(h)
	defer s.Close()

	err := requests.URL(s.URL + "/api/stock-prices?stock=TSLA&like=true").Fetch(context.Background())
	require.NoError(t, err)
	update := bson.D{{"$inc", bson.D{{"likes_count", 2}}}}
	_, err = stockServ.stocks.UpdateByID(context.Background(), "TSLA", update)
	require.NoError(t, err)
	err = requests.URL(s.URL + "/api/stock-prices?stock=KO&like=true").Fetch(context.Background())
	require.NoError(t, err)

	var actual StockPricesResp
	err = requests.URL(s.URL + "/api/stock-prices?stock=TSLA&stock=KO&like=true").ToJSON(&actual).Fetch(context.Background())

	require.NoError(t, err)
	require.Len(t, actual.StockData, 2)
	assert.NotZero(t, actual.StockData[0].Price)
	assert.NotZero(t, actual.StockData[1].Price)

	expected := StockPricesResp{
		StockData: []StockDatasResp{
			{
				Stock:    "TSLA",
				Price:    actual.StockData[0].Price,
				RelLikes: 2,
			},
			{
				Stock:    "KO",
				Price:    actual.StockData[1].Price,
				RelLikes: -2,
			},
		},
	}
	assert.Equal(t, expected, actual)
}

func newTestServer() *httptest.Server {
	h := NewHandler(echo.New(), NewStockService(newTestMongoDatabase()))

	return httptest.NewServer(h)
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
