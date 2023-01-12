package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	chi "github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"messageboard/api"
	"messageboard/httpserv"
	"messageboard/thread"
)

var (
	client = http.Client{Timeout: 2 * time.Second}
	db     = newTestMongoDatabase()
)

func init() {
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

	var res *http.Response
	var err error
	ifForm := rand.Int()%2 == 0
	if ifForm {
		t.Log("post form data")
		res, err = client.PostForm(s.URL+"/api/threads/board_test", url.Values{
			"text":            {"Some text."},
			"delete_password": {"p@ssw0rd"},
		})
	} else {
		t.Log("post json data")
		body, _ := json.Marshal(api.CreateThreadBody{
			DeletePassword: "p@ssw0rd",
			Text:           "Some text.",
		})
		res, err = client.Post(s.URL+"/api/threads/board_test", "application/json", bytes.NewReader(body))
	}
	if err != nil {
		t.Fatal(err)
	}

	if http.StatusOK != res.StatusCode {
		t.Fatalf("expected '200 OK' status, got '%s'", res.Status)
	}

	res, err = client.Get(s.URL + "/api/threads/board_test")
	if err != nil {
		t.Fatal(err)
	}

	if http.StatusOK != res.StatusCode {
		t.Fatalf("expected '200 OK' status, got '%s'", res.Status)
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	_ = res.Body.Close()
	res.Body = io.NopCloser(bytes.NewBuffer(resBytes))
	actual := string(resBytes)

	var threads []api.Thread
	if err := json.NewDecoder(res.Body).Decode(&threads); err != nil {
		t.Fatal(err)
	}
	_ = res.Body.Close()

	createdOn := threads[0].CreatedOn.Format(time.RFC3339)
	bumpedOn := threads[0].BumpedOn.Format(time.RFC3339)
	expected := fmt.Sprintf(`[{"_id":"%s","bumped_on":"%s","created_on":"%s","replies":[],"text":"Some text."}]
`, threads[0].Id, bumpedOn, createdOn)
	if expected != actual {
		t.Fatalf("expected %+v, got %+v", expected, actual)
	}

	resBytes, err = io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
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

	res := mongoClient.Database("test_message_board")

	_, _ = res.Collection(thread.ThreadsCollection).DeleteMany(context.Background(), bson.D{{}})

	return res
}
