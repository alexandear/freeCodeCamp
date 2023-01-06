package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client = &http.Client{
	Timeout: time.Second,
}

func TestHandlerServeStaticContent(t *testing.T) {
	h := newHandler(echo.New(), nil)

	s := httptest.NewServer(h)
	defer s.Close()

	resIndex, err := client.Get(s.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resIndex.Body.Close()

	indexContent, err := io.ReadAll(resIndex.Body)
	if !strings.Contains(string(indexContent), "<h1>Exercise tracker</h1>") {
		t.Fatal("index.html must be served")
	}

	resStyle, err := client.Get(s.URL + "/style.css")
	if err != nil {
		t.Fatal(err)
	}

	styleContent, err := io.ReadAll(resStyle.Body)
	if !strings.Contains(string(styleContent), "background-color:") {
		t.Fatal("style.css must be served")
	}
}

func TestHandler_CreateUser(t *testing.T) {
	mg := newTestMongoClient(t)
	h := newHandler(echo.New(), mg)

	s := httptest.NewServer(h)
	defer s.Close()

	expected := user{
		Username: "johndoe",
	}
	res, err := client.PostForm(s.URL+"/api/users", url.Values{
		"username": {expected.Username},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	var actual user
	if err := json.NewDecoder(res.Body).Decode(&actual); err != nil {
		t.Fatal(err)
	}

	if actual.ID == "" {
		t.Fatalf("_id must be non-empty")
	}
	expected.ID = actual.ID

	if expected != actual {
		t.Fatalf("expected %+v, got %+v", expected, actual)
	}
}

func newTestMongoClient(t *testing.T) *mongo.Client {
	t.Helper()

	mongoURI := os.Getenv("MONGODB_URI")

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPIOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		t.Fatal(err)
	}

	return mongoClient
}
