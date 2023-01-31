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
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestHandler_StockPrice_GetOneStock(t *testing.T) {
	t.Parallel()
	s := newTestServer(t)
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
	t.Parallel()
	s := newTestServer(t)
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
	t.Parallel()
	s := newTestServer(t)
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
	t.Parallel()
	s := newTestServer(t)
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
	t.Parallel()
	s := newTestServer(t)
	defer s.Close()

	reqLike := requests.URL(s.URL + "/api/stock-prices?stock=TSLA&like=true")
	require.NoError(t, reqLike.Header("X-Real-Ip", "172.27.0.31").Fetch(context.Background()))
	require.NoError(t, reqLike.Header("X-Real-Ip", "172.27.0.32").Fetch(context.Background()))
	err := requests.URL(s.URL + "/api/stock-prices?stock=KO&like=true").Fetch(context.Background())
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

type testServer struct {
	URL string

	t *testing.T

	server      *httptest.Server
	mongoClient *mongo.Client
}

func newTestServer(t *testing.T) *testServer {
	t.Helper()

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		t.Skip("MONGODB_URI not set, skipping test")
	}

	var client *mongo.Client
	func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
		clientOptions := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPIOptions)
		cl, err := mongo.Connect(ctx, clientOptions)
		require.NoError(t, err)
		client = cl
	}()

	db := client.Database("test_stock_checker")

	stockServer := NewStockService(db)
	func() {
		ctxDel, cancelDel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelDel()

		_, err := stockServer.stocks.DeleteMany(ctxDel, bson.D{{}})
		require.NoError(t, err)
		_, err = stockServer.stockPerIPs.DeleteMany(ctxDel, bson.D{{}})
		require.NoError(t, err)
	}()

	h := NewHandler(echo.New(), stockServer)
	ts := httptest.NewServer(h)
	return &testServer{
		URL:         ts.URL,
		t:           t,
		server:      ts,
		mongoClient: client,
	}
}

func (ts *testServer) Close() {
	ts.t.Helper()
	ts.server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	require.NoError(ts.t, ts.mongoClient.Disconnect(ctx))
	cancel()
}
