package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
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

	res, err := client.PostForm(s.URL+"/api/users", url.Values{
		"username": {"johndoe"},
	})
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
	var u handlerUser
	if err := json.NewDecoder(bytes.NewReader(resBytes)).Decode(&u); err != nil {
		t.Fatal(err)
	}

	if u.ID == "" {
		t.Fatalf("_id must be non empty")
	}

	expected := fmt.Sprintf(`{"_id":"%s","username":"johndoe"}
`, u.ID)
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

	if http.StatusOK != res.StatusCode {
		t.Fatalf("expected '200 OK' status, got '%s'", res.Status)
	}

	var actual []handlerUser
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
	var u handlerUser
	if err := json.NewDecoder(resUser.Body).Decode(&u); err != nil {
		t.Fatal(err)
	}

	res, err := client.PostForm(s.URL+"/api/users/"+u.ID+"/exercises", url.Values{
		"description": {"test"},
		"duration":    {"60"},
		"date":        {"1990-01-01"},
	})
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

	expected := fmt.Sprintf(`{"_id":"%s","username":"johndoe","date":"Mon Jan 01 1990","duration":60,"description":"test"}
`, u.ID)
	if expected != actual {
		t.Fatalf("expected %+v, got %+v", expected, actual)
	}
}

func TestHandler_Logs(t *testing.T) {
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
	var u handlerUser
	if err := json.NewDecoder(resUser.Body).Decode(&u); err != nil {
		t.Fatal(err)
	}

	for _, ex := range []handlerExercise{
		{
			Description: "ex 1",
			Duration:    30,
			Date:        "2023-02-22",
		},
		{
			Description: "ex 2",
			Duration:    45,
			Date:        "2023-02-25",
		},
	} {
		if _, err := client.PostForm(s.URL+"/api/users/"+u.ID+"/exercises", url.Values{
			"description": {ex.Description},
			"duration":    {strconv.Itoa(ex.Duration)},
			"date":        {ex.Date},
		}); err != nil {
			t.Fatal(err)
		}
	}

	res, err := client.Get(s.URL + "/api/users/" + u.ID + "/logs")
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

	expected := fmt.Sprintf(`{"_id":"%s","username":"johndoe","count":2,"log":[{"date":"Wed Feb 02 2023","duration":30,"description":"ex 1"},{"date":"Sat Feb 02 2023","duration":45,"description":"ex 2"}]}
`, u.ID)
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
