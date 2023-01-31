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

func TestHandlerServeStaticContent(t *testing.T) {
	t.Parallel()
	h := newHandler(echo.New(), nil)

	s := httptest.NewServer(h)
	defer s.Close()

	resIndex, err := http.Get(s.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resIndex.Body.Close()

	indexContent, err := io.ReadAll(resIndex.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(indexContent), "<h1>Exercise tracker</h1>") {
		t.Fatal("index.html must be served")
	}

	resStyle, err := http.Get(s.URL + "/style.css")
	if err != nil {
		t.Fatal(err)
	}

	styleContent, err := io.ReadAll(resStyle.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(styleContent), "background-color:") {
		t.Fatal("style.css must be served")
	}
}

func TestHandler_CreateUser(t *testing.T) {
	t.Parallel()

	s := newTestServer(t)
	defer s.Close()

	res, err := http.PostForm(s.URL+"/api/users", url.Values{
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
	t.Parallel()

	s := newTestServer(t)
	defer s.Close()

	for _, username := range []string{"cat", "wolf"} {
		if _, err := http.PostForm(s.URL+"/api/users", url.Values{
			"username": {username},
		}); err != nil {
			t.Fatal(err)
		}
	}

	res, err := http.Get(s.URL + "/api/users")
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
	t.Parallel()

	s := newTestServer(t)
	defer s.Close()

	resUser, err := http.PostForm(s.URL+"/api/users", url.Values{
		"username": {"johndoe"},
	})
	if err != nil {
		t.Fatal(err)
	}
	var u handlerUser
	if err := json.NewDecoder(resUser.Body).Decode(&u); err != nil {
		t.Fatal(err)
	}

	res, err := http.PostForm(s.URL+"/api/users/"+u.ID+"/exercises", url.Values{
		"description": {"test"},
		"duration":    {"60"},
		"date":        {"1990-01-10"},
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

	expected := fmt.Sprintf(`{"_id":"%s","username":"johndoe","date":"Wed Jan 10 1990","duration":60,"description":"test"}
`, u.ID)
	if expected != actual {
		t.Fatalf("expected %+v, got %+v", expected, actual)
	}
}

func TestHandler_Logs(t *testing.T) {
	t.Parallel()

	s := newTestServer(t)
	defer s.Close()

	resUser, err := http.PostForm(s.URL+"/api/users", url.Values{
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
			Date:        "2022-10-22",
		},
		{
			Description: "ex 2",
			Duration:    45,
			Date:        "2023-02-09",
		},
		{
			Description: "ex 3",
			Duration:    10,
			Date:        "2023-11-25",
		},
	} {
		if _, err := http.PostForm(s.URL+"/api/users/"+u.ID+"/exercises", url.Values{
			"description": {ex.Description},
			"duration":    {strconv.Itoa(ex.Duration)},
			"date":        {ex.Date},
		}); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("all", func(t *testing.T) {
		res, err := http.Get(s.URL + "/api/users/" + u.ID + "/logs")
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

		expected := fmt.Sprintf(`{"_id":"%s","username":"johndoe","count":3,"log":[{"date":"Sat Oct 22 2022","duration":30,"description":"ex 1"},{"date":"Thu Feb 09 2023","duration":45,"description":"ex 2"},{"date":"Sat Nov 25 2023","duration":10,"description":"ex 3"}]}
`, u.ID)
		if expected != actual {
			t.Fatalf("expected %+v, got %+v", expected, actual)
		}
	})

	t.Run("limit", func(t *testing.T) {
		res, err := http.Get(s.URL + "/api/users/" + u.ID + "/logs?limit=1")
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

		expected := fmt.Sprintf(`{"_id":"%s","username":"johndoe","count":1,"log":[{"date":"Sat Oct 22 2022","duration":30,"description":"ex 1"}]}
`, u.ID)
		if expected != actual {
			t.Fatalf("expected %+v, got %+v", expected, actual)
		}
	})

	t.Run("from", func(t *testing.T) {
		res, err := http.Get(s.URL + "/api/users/" + u.ID + "/logs?from=2023-01-01")
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

		expected := fmt.Sprintf(`{"_id":"%s","username":"johndoe","count":2,"log":[{"date":"Thu Feb 09 2023","duration":45,"description":"ex 2"},{"date":"Sat Nov 25 2023","duration":10,"description":"ex 3"}]}
`, u.ID)
		if expected != actual {
			t.Fatalf("expected %+v, got %+v", expected, actual)
		}
	})

	t.Run("to", func(t *testing.T) {
		res, err := http.Get(s.URL + "/api/users/" + u.ID + "/logs?to=2023-04-01")
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

		expected := fmt.Sprintf(`{"_id":"%s","username":"johndoe","count":2,"log":[{"date":"Sat Oct 22 2022","duration":30,"description":"ex 1"},{"date":"Thu Feb 09 2023","duration":45,"description":"ex 2"}]}
`, u.ID)
		if expected != actual {
			t.Fatalf("expected %+v, got %+v", expected, actual)
		}
	})
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
		if err != nil {
			t.Fatal(err)
		}
		client = cl
	}()

	db := client.Database("test_exercise_tracker")

	us := newUserService(db)
	h := newHandler(echo.New(), us)
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
	defer cancel()
	if err := ts.mongoClient.Disconnect(ctx); err != nil {
		ts.t.Fatal(err)
	}
}
