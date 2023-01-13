package test

import (
	"context"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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
	"messageboard/msgboard"
	clapi "messageboard/test/client/api"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 --config=client/client.cfg.yaml ../api/openapi.yaml

func init() {
	gofakeit.Seed(0)
	rand.Seed(time.Now().UnixNano())
}

func TestCreateNewThread(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	client := newTestClient(t, s.URL)

	board := gofakeit.Animal()
	text := gofakeit.BuzzWord()

	var createResp *clapi.CreateThreadResponse
	createBody := clapi.CreateThreadJSONRequestBody{
		Text:           text,
		DeletePassword: gofakeit.NounAbstract(),
	}
	if rand.Int()%2 == 0 {
		var err error
		createResp, err = client.CreateThreadWithFormdataBodyWithResponse(context.Background(), board, createBody)
		require.NoError(t, err)
	} else {
		var err error
		createResp, err = client.CreateThreadWithResponse(context.Background(), board, createBody)
		require.NoError(t, err)
	}
	assert.Equal(t, http.StatusFound, createResp.StatusCode())

	threadID := threadIDFromHeader(createResp.HTTPResponse.Header)
	assert.NotEmpty(t, threadID)

	getResp, err := client.GetThreadsWithResponse(context.Background(), board)
	require.NoError(t, err)

	threads := *getResp.JSON200
	require.Len(t, threads, 1)
	thread := threads[0]
	assert.Equal(t, threadID, thread.Id)
	assert.NotZero(t, thread.CreatedOn)
	assert.Equal(t, thread.CreatedOn, thread.BumpedOn)
	assert.Equal(t, text, thread.Text)
	assert.Len(t, thread.Replies, 0)
	assert.Equal(t, 0, thread.Replycount)
}

func TestViewTheMost10RecentThreadsWith3RepliesEach(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	client := newTestClient(t, s.URL)

	board := gofakeit.Animal()

	for thread := 0; thread < 15; thread++ {
		resp, err := client.CreateThreadWithResponse(context.Background(), board, clapi.CreateThreadBody{
			DeletePassword: gofakeit.NounAbstract(),
			Text:           gofakeit.BuzzWord(),
		})
		require.NoError(t, err)

		threadID := threadIDFromHeader(resp.HTTPResponse.Header)

		for reply := 0; reply < 5; reply++ {
			_, err = client.CreateReplyWithResponse(context.Background(), board, clapi.CreateReplyBody{
				DeletePassword: gofakeit.NounAbstract(),
				Text:           gofakeit.BuzzWord(),
				ThreadId:       threadID,
			})
			require.NoError(t, err)
		}
	}

	getResp, err := client.GetThreadsWithResponse(context.Background(), board)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, getResp.StatusCode())

	threads := *getResp.JSON200
	assert.Len(t, threads, 10)
	for _, thread := range threads {
		assert.Len(t, thread.Replies, 3)
	}
}

func TestDeleteThreadWithIncorrectPassword(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	client := newTestClient(t, s.URL)

	board := gofakeit.Animal()

	createResp, err := client.CreateThreadWithResponse(context.Background(), board, clapi.CreateThreadJSONRequestBody{
		Text:           gofakeit.BuzzWord(),
		DeletePassword: gofakeit.NounAbstract(),
	})
	require.NoError(t, err)

	threadID := threadIDFromHeader(createResp.HTTPResponse.Header)

	deleteResp, err := client.DeleteThreadWithResponse(context.Background(), board, clapi.DeleteThreadJSONRequestBody{
		DeletePassword: "wrong password",
		ThreadId:       threadID,
	})
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, deleteResp.StatusCode())
	assert.Equal(t, "incorrect password", string(deleteResp.Body))
}

func TestDeleteThreadWithCorrectPassword(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	client := newTestClient(t, s.URL)

	board := gofakeit.Animal()
	deletePassword := gofakeit.NounAbstract()

	createResp, err := client.CreateThreadWithResponse(context.Background(), board, clapi.CreateThreadJSONRequestBody{
		Text:           gofakeit.BuzzWord(),
		DeletePassword: deletePassword,
	})
	require.NoError(t, err)

	threadID := threadIDFromHeader(createResp.HTTPResponse.Header)

	deleteResp, err := client.DeleteThreadWithResponse(context.Background(), board, clapi.DeleteThreadJSONRequestBody{
		DeletePassword: deletePassword,
		ThreadId:       threadID,
	})
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, deleteResp.StatusCode())
	assert.Equal(t, "success", string(deleteResp.Body))
}

func TestReportThread(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	client := newTestClient(t, s.URL)

	board := gofakeit.Animal()
	text := gofakeit.BuzzWord()
	deletePassword := gofakeit.NounAbstract()

	createResp, err := client.CreateThreadWithResponse(context.Background(), board, clapi.CreateThreadJSONRequestBody{
		Text:           text,
		DeletePassword: deletePassword,
	})
	require.NoError(t, err)

	threadID := threadIDFromHeader(createResp.HTTPResponse.Header)

	reportResp, err := client.ReportThreadWithResponse(context.Background(), board, clapi.ReportThreadBody{
		ThreadId: threadID,
	})
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, reportResp.StatusCode())
	assert.Equal(t, "reported", string(reportResp.Body))
}

func TestCreateNewReply(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	client := newTestClient(t, s.URL)

	board := gofakeit.Animal()
	threadText := gofakeit.BuzzWord()

	threadID := func() string {
		resp, err := client.CreateThreadWithResponse(context.Background(), board, clapi.CreateThreadBody{
			DeletePassword: gofakeit.NounAbstract(),
			Text:           threadText,
		})
		require.NoError(t, err)
		return threadIDFromHeader(resp.HTTPResponse.Header)
	}()

	replyText := gofakeit.BuzzWord()

	var createResp *clapi.CreateReplyResponse
	createBody := clapi.CreateReplyBody{
		DeletePassword: gofakeit.NounAbstract(),
		Text:           replyText,
		ThreadId:       threadID,
	}
	if rand.Int()%2 == 0 {
		var err error
		createResp, err = client.CreateReplyWithFormdataBodyWithResponse(context.Background(), board, createBody)
		require.NoError(t, err)
	} else {
		var err error
		createResp, err = client.CreateReplyWithResponse(context.Background(), board, createBody)
		require.NoError(t, err)
	}

	assert.Equal(t, http.StatusOK, createResp.StatusCode())
	replyID := replyIDFromHeader(createResp.HTTPResponse.Header)
	assert.NotEmpty(t, replyID)

	getResp, err := client.GetRepliesWithResponse(context.Background(), board, &clapi.GetRepliesParams{ThreadId: threadID})
	require.NoError(t, err)
	thread := *getResp.JSON200
	assert.True(t, thread.BumpedOn.After(thread.CreatedOn))
	assert.Len(t, thread.Replies, 1)
	assert.Equal(t, 1, thread.Replycount)
	reply := thread.Replies[0]
	assert.Equal(t, replyID, reply.Id)
	assert.Equal(t, replyText, reply.Text)
	assert.Equal(t, thread.BumpedOn, reply.CreatedOn)
}

func TestViewThreadWithAllReplies(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	client := newTestClient(t, s.URL)

	board := gofakeit.Animal()
	threadText := gofakeit.BuzzWord()

	threadID := func() string {
		resp, err := client.CreateThreadWithResponse(context.Background(), board, clapi.CreateThreadBody{
			DeletePassword: gofakeit.NounAbstract(),
			Text:           threadText,
		})
		require.NoError(t, err)
		return threadIDFromHeader(resp.HTTPResponse.Header)
	}()

	for i := 0; i < 5; i++ {
		_, err := client.CreateReplyWithResponse(context.Background(), board, clapi.CreateReplyBody{
			DeletePassword: gofakeit.NounAbstract(),
			Text:           gofakeit.BuzzWord(),
			ThreadId:       threadID,
		})
		require.NoError(t, err)
	}

	getResp, err := client.GetRepliesWithResponse(context.Background(), board, &clapi.GetRepliesParams{
		ThreadId: threadID,
	})
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, getResp.StatusCode())
	thread := getResp.JSON200
	assert.Equal(t, threadText, thread.Text)
	assert.Len(t, thread.Replies, 5)
}

func TestDeleteReplyWithIncorrectPassword(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	client := newTestClient(t, s.URL)

	board := gofakeit.Animal()

	threadID := func() string {
		resp, err := client.CreateThreadWithResponse(context.Background(), board, clapi.CreateThreadBody{
			DeletePassword: gofakeit.NounAbstract(),
			Text:           gofakeit.BuzzWord(),
		})
		require.NoError(t, err)
		return threadIDFromHeader(resp.HTTPResponse.Header)
	}()

	replyText := gofakeit.BuzzWord()

	deletePassword := gofakeit.NounAbstract()
	createBody := clapi.CreateReplyBody{
		DeletePassword: deletePassword,
		Text:           replyText,
		ThreadId:       threadID,
	}
	createResp, err := client.CreateReplyWithResponse(context.Background(), board, createBody)
	require.NoError(t, err)

	replyID := replyIDFromHeader(createResp.HTTPResponse.Header)

	deleteReply, err := client.DeleteReplyWithResponse(context.Background(), board, clapi.DeleteReplyBody{
		DeletePassword: "incorrect password",
		ReplyId:        replyID,
		ThreadId:       threadID,
	})
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, deleteReply.StatusCode())
	assert.Equal(t, "incorrect password", string(deleteReply.Body))
}

func TestDeleteReplyWithCorrectPassword(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	client := newTestClient(t, s.URL)

	board := gofakeit.Animal()
	threadText := gofakeit.BuzzWord()

	threadID := func() string {
		resp, err := client.CreateThreadWithResponse(context.Background(), board, clapi.CreateThreadBody{
			DeletePassword: gofakeit.NounAbstract(),
			Text:           threadText,
		})
		require.NoError(t, err)
		return threadIDFromHeader(resp.HTTPResponse.Header)
	}()

	replyText := gofakeit.BuzzWord()

	deletePassword := gofakeit.NounAbstract()
	createBody := clapi.CreateReplyBody{
		DeletePassword: deletePassword,
		Text:           replyText,
		ThreadId:       threadID,
	}
	createResp, err := client.CreateReplyWithResponse(context.Background(), board, createBody)
	require.NoError(t, err)

	replyID := replyIDFromHeader(createResp.HTTPResponse.Header)

	deleteReply, err := client.DeleteReplyWithResponse(context.Background(), board, clapi.DeleteReplyBody{
		DeletePassword: deletePassword,
		ReplyId:        replyID,
		ThreadId:       threadID,
	})
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, deleteReply.StatusCode())
	assert.Equal(t, "success", string(deleteReply.Body))

	getReply, err := client.GetRepliesWithResponse(context.Background(), board, &clapi.GetRepliesParams{ThreadId: threadID})
	require.NoError(t, err)

	thread := *getReply.JSON200
	reply := thread.Replies[0]
	assert.Equal(t, replyID, reply.Id)
	assert.Equal(t, "[deleted]", reply.Text)
}

func TestReportReply(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	client := newTestClient(t, s.URL)

	board := gofakeit.Animal()
	threadText := gofakeit.BuzzWord()

	threadID := func() string {
		resp, err := client.CreateThreadWithResponse(context.Background(), board, clapi.CreateThreadBody{
			DeletePassword: gofakeit.NounAbstract(),
			Text:           threadText,
		})
		require.NoError(t, err)
		return threadIDFromHeader(resp.HTTPResponse.Header)
	}()

	replyText := gofakeit.BuzzWord()

	deletePassword := gofakeit.NounAbstract()
	createBody := clapi.CreateReplyBody{
		DeletePassword: deletePassword,
		Text:           replyText,
		ThreadId:       threadID,
	}
	createResp, err := client.CreateReplyWithResponse(context.Background(), board, createBody)
	require.NoError(t, err)

	replyID := replyIDFromHeader(createResp.HTTPResponse.Header)

	reportResp, err := client.ReportReplyWithResponse(context.Background(), board, clapi.ReportReplyBody{
		ReplyId:  replyID,
		ThreadId: threadID,
	})
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, reportResp.StatusCode())
	assert.Equal(t, "reported", string(reportResp.Body))
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

	_, _ = res.Collection(msgboard.ThreadsCollection).DeleteMany(context.Background(), bson.D{{}})
	_, _ = res.Collection(msgboard.RepliesCollection).DeleteMany(context.Background(), bson.D{{}})

	return res
}

func newTestServer() *httptest.Server {
	db := newTestMongoDatabase()
	msgServ := msgboard.NewService(db)
	serv := httpserv.NewServer(msgServ)
	r := chi.NewRouter()
	strictHandler := api.NewStrictHandler(serv, nil)
	api.HandlerFromMux(strictHandler, r)
	return httptest.NewServer(r)
}

func newTestClient(t *testing.T, serverURL string) *clapi.ClientWithResponses {
	client, err := clapi.NewClientWithResponses(serverURL, clapi.WithHTTPClient(&http.Client{
		Timeout: 2 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}))
	require.NoError(t, err)
	return client
}

func threadIDFromHeader(r http.Header) string {
	loc := r.Get("Location")
	paths := strings.Split(loc, "/")
	if len(paths) != 4 {
		return ""
	}
	return paths[3]
}

func replyIDFromHeader(r http.Header) string {
	return r.Get("X-Message-Board-Reply-ID")
}
