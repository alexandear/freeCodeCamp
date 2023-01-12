package test

import (
	"context"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	chi "github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"messageboard/api"
	"messageboard/httpserv"
	clapi "messageboard/test/client/api"
	"messageboard/thread"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 --config=client/client.cfg.yaml ../api/openapi.yaml

var (
	db = newTestMongoDatabase()
)

func init() {
	gofakeit.Seed(0)
	rand.Seed(time.Now().UnixNano())
}

func TestCreateNewThread(t *testing.T) {
	threadServ := thread.NewService(db)
	serv := httpserv.NewServer(threadServ)
	r := chi.NewRouter()
	strictHandler := api.NewStrictHandler(serv, nil)
	api.HandlerFromMux(strictHandler, r)
	s := httptest.NewServer(r)
	defer s.Close()

	client, err := clapi.NewClientWithResponses(s.URL, clapi.WithHTTPClient(&http.Client{Timeout: 2 * time.Second}))
	require.NoError(t, err)

	board := gofakeit.Animal()
	text := gofakeit.BuzzWord()
	deletePassword := gofakeit.NounAbstract()

	var createResp *clapi.CreateThreadResponse
	createBody := clapi.CreateThreadJSONRequestBody{
		Text:           text,
		DeletePassword: deletePassword,
	}
	if rand.Int()%2 == 0 {
		createResp, err = client.CreateThreadWithFormdataBodyWithResponse(context.Background(), board, createBody)
		require.NoError(t, err)
	} else {
		createResp, err = client.CreateThreadWithResponse(context.Background(), board, createBody)
		require.NoError(t, err)
	}
	assert.Equal(t, createResp.StatusCode(), http.StatusOK)

	getResp, err := client.GetThreadsWithResponse(context.Background(), board)
	require.NoError(t, err)

	threads := *getResp.JSON200
	require.Len(t, threads, 1)
	thread := threads[0]
	assert.NotEmpty(t, thread.Id)
	assert.NotZero(t, thread.CreatedOn)
	assert.Equal(t, thread.CreatedOn, thread.BumpedOn)
	assert.Equal(t, text, thread.Text)
	assert.Len(t, thread.Replies, 0)
}

func TestCreateReply(t *testing.T) {
	threadServ := thread.NewService(db)
	serv := httpserv.NewServer(threadServ)
	r := chi.NewRouter()
	strictHandler := api.NewStrictHandler(serv, nil)
	api.HandlerFromMux(strictHandler, r)
	s := httptest.NewServer(r)
	defer s.Close()

	client, err := clapi.NewClientWithResponses(s.URL, clapi.WithHTTPClient(&http.Client{Timeout: 2 * time.Second}))
	require.NoError(t, err)

	board := gofakeit.Animal()
	threadText := gofakeit.BuzzWord()

	threadID := func() string {
		_, err := client.CreateThreadWithResponse(context.Background(), board, clapi.CreateThreadBody{
			DeletePassword: gofakeit.NounAbstract(),
			Text:           threadText,
		})
		require.NoError(t, err)

		getResp, err := client.GetThreadsWithResponse(context.Background(), board)
		require.NoError(t, err)
		require.Len(t, *getResp.JSON200, 1)
		return (*getResp.JSON200)[0].Id
	}()

	replyText := gofakeit.BuzzWord()

	var createResp *clapi.CreateReplyResponse
	createBody := clapi.CreateReplyBody{
		DeletePassword: gofakeit.NounAbstract(),
		Text:           replyText,
		ThreadId:       threadID,
	}
	if rand.Int()%2 == 0 {
		createResp, err = client.CreateReplyWithFormdataBodyWithResponse(context.Background(), board, createBody)
		require.NoError(t, err)
	} else {
		createResp, err = client.CreateReplyWithResponse(context.Background(), board, createBody)
		require.NoError(t, err)
	}

	assert.Equal(t, createResp.StatusCode(), http.StatusOK)

	getResp, err := client.GetThreadsWithResponse(context.Background(), board)
	require.NoError(t, err)
	require.Len(t, *getResp.JSON200, 1)
	thread := (*getResp.JSON200)[0]
	assert.True(t, thread.BumpedOn.After(thread.CreatedOn))
	assert.Len(t, thread.Replies, 1)
	reply := thread.Replies[0]
	assert.NotEmpty(t, reply.Id)
	assert.Equal(t, replyText, reply.Text)
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

	res := mongoClient.Database("test_message_board")

	_, _ = res.Collection(thread.ThreadsCollection).DeleteMany(context.Background(), bson.D{{}})
	_, _ = res.Collection(thread.RepliesCollection).DeleteMany(context.Background(), bson.D{{}})

	return res
}
