package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockClock struct{}

func (c *mockClock) Now() time.Time {
	return time.Date(2023, time.January, 10, 23, 0, 2, 0, time.UTC)
}

func TestApiHandler_OK(t *testing.T) {
	t.Parallel()
	api := &apiHandler{
		clock: &mockClock{},
	}

	s := httptest.NewServer(api)
	defer s.Close()

	for name, tc := range map[string]struct {
		url  string
		resp respSuccess
	}{
		"valid date": {
			url: "/api/2016-12-25",
			resp: respSuccess{
				UnixMs: 1482624000000,
				UTC:    "Sun, 25 Dec 2016 00:00:00 GMT",
			},
		},
		"valid unix ms": {
			url: "/api/1451001600000",
			resp: respSuccess{
				UnixMs: 1451001600000,
				UTC:    "Fri, 25 Dec 2015 02:00:00 GMT",
			},
		},
		"custom date format": {
			url: "/api/05 October 2011, GMT",
			resp: respSuccess{
				UnixMs: 1317772800000,
				UTC:    "Wed, 05 Oct 2011 00:00:00 GMT",
			},
		},
		"empty date parameter": {
			url: "/api",
			resp: respSuccess{
				UnixMs: 1673391602000,
				UTC:    "Tue, 10 Jan 2023 23:00:02 GMT",
			},
		},
		"empty date parameter with slash": {
			url: "/api/",
			resp: respSuccess{
				UnixMs: 1673391602000,
				UTC:    "Tue, 10 Jan 2023 23:00:02 GMT",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			res, err := http.Get(s.URL + tc.url)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			if http.StatusOK != res.StatusCode {
				t.Fatal("response must be OK")
			}

			var resp respSuccess
			if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
				t.Fatal(err)
			}
			if tc.resp != resp {
				t.Fatalf("expected %+v, got %+v", tc.resp, resp)
			}
		})
	}
}

func TestApiHandler_Error(t *testing.T) {
	t.Parallel()
	api := &apiHandler{}

	s := httptest.NewServer(api)
	defer s.Close()

	req, err := http.Get(s.URL + "/api/this-is-not-a-date")
	if err != nil {
		t.Fatal(err)
	}

	if req.StatusCode != http.StatusBadRequest {
		t.Fatal("response must be Bad request")
	}

	var resp respError
	if err := json.NewDecoder(req.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	expected := respError{Error: "Invalid Date"}
	if expected != resp {
		t.Fatalf("expected %+v, got %+v", expected, resp)
	}
}
