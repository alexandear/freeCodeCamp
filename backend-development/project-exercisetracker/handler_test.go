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

var db *mongo.Database

func TestMain(m *testing.M) {
	db = newTestMongoDatabase()
	exit := m.Run()
	os.Exit(exit)
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
	us := newUserService(db)
	h := newHandler(echo.New(), us)

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

func TestHandler_Users(t *testing.T) {
	us := newUserService(db)
	h := newHandler(echo.New(), us)

	s := httptest.NewServer(h)
	defer s.Close()

	for _, username := range []string{"cat", "wolf"} {
		if _, err := client.PostForm(s.URL+"/api/users", url.Values{
			"username": {username},
		}); err != nil {
			t.Fatal(err)
		}
	}

	res, err := client.Get(s.URL + "/api/users")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	var actual []user
	if err := json.NewDecoder(res.Body).Decode(&actual); err != nil {
		t.Fatal(err)
	}

	if len(actual) == 0 {
		t.Fatal("users must be returned")
	}

	for _, u := range actual {
		if u.ID == "" || u.Username == "" {
			t.Fatalf("id and username must not be empty, got: %+v", u)
		}
	}
}

func TestHandler_CreateExercise(t *testing.T) {
	us := newUserService(db)
	h := newHandler(echo.New(), us)

	s := httptest.NewServer(h)
	defer s.Close()

	resUser, err := client.PostForm(s.URL+"/api/users", url.Values{
		"username": {"johndoe"},
	})
	if err != nil {
		t.Fatal(err)
	}
	var u user
	if err := json.NewDecoder(resUser.Body).Decode(&u); err != nil {
		t.Fatal(err)
	}

	expected := exercise{
		Description: "test",
		Duration:    60,
		Date:        "Mon Jan 01 1990",
		User: user{
			ID:       u.ID,
			Username: "johndoe",
		},
	}
	res, err := client.PostForm(s.URL+"/api/users/"+u.ID+"/exercises", url.Values{
		"description": {"test"},
		"duration":    {"60"},
		"date":        {"1990-01-01"},
	})
	if err != nil {
		t.Fatal(err)
	}

	if http.StatusOK != res.StatusCode {
		t.Fatal("status must be OK")
	}

	var actual exercise
	if err := json.NewDecoder(res.Body).Decode(&actual); err != nil {
		t.Fatal(err)
	}

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

	return mongoClient.Database("test_exercise_tracker")
}
