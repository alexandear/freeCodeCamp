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
	api := &apiHandler{
		clock: &mockClock{},
	}

	s := httptest.NewServer(api)
	defer s.Close()

	for _, tc := range []struct {
		url       string
		expUnixMs *int64
		expUTC    *string
	}{
		{
			url:       "/api",
			expUnixMs: int64Ptr(1673391602000),
			expUTC:    stringPtr("Fri, 10 Jan 2023 23:00:02 GMT"),
		},
		{
			url:       "/api/",
			expUnixMs: int64Ptr(1673391602000),
			expUTC:    stringPtr("Fri, 10 Jan 2023 23:00:02 GMT"),
		},
		// {
		// 	url:       "/api/2016-12-25",
		// 	expUnixMs: int64Ptr(1673391602000),
		// 	expUTC:    stringPtr("Fri, 10 Jan 2023 23:00:02 GMT"),
		// },
		// {
		// 	url:       "/api/1451001600000",
		// 	expUnixMs: int64Ptr(1673391602000),
		// 	expUTC:    stringPtr("Fri, 10 Jan 2023 23:00:02 GMT"),
		// },
		// {
		// 	url:       "/api/05 October 2011, GMT",
		// 	expUnixMs: int64Ptr(1673391602000),
		// 	expUTC:    stringPtr("Fri, 10 Jan 2023 23:00:02 GMT"),
		// },
	} {
		t.Run(tc.url, func(t *testing.T) {
			req, err := http.Get(s.URL + tc.url)
			if err != nil {
				t.FailNow()
			}

			resp := &respSuccess{}
			if err := json.NewDecoder(req.Body).Decode(resp); err != nil {
				t.FailNow()
			}
			if tc.expUnixMs != nil && *tc.expUnixMs != *resp.UnixMs {
				t.Fatalf("expected %d, got %d", *tc.expUnixMs, *resp.UnixMs)
			}
      if tc.expUTC != nil && *tc.expUTC != resp.UTC {
				t.Fatalf("expected %s, got %s", *tc.expUTC, resp.UTC)
			}
		})
	}
}

func TestApiHandler_Error(t *testing.T) {
	api := &apiHandler{}

	s := httptest.NewServer(api)
	defer s.Close()

  req, err := http.Get(s.URL + "/api/this-is-not-a-date")
  if err != nil {
    t.FailNow()
  }

  resp := &respError{}
  if err := json.NewDecoder(req.Body).Decode(resp); err != nil {
    t.FailNow()
  }
  if "Invalid Date" != resp.Error {
    t.FailNow()
  }

}

func int64Ptr(val int64) *int64 {
	return &val
}

func stringPtr(val string) *string {
	return &val
}
