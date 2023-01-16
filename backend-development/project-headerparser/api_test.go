package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/imroc/req/v3"
)

func TestApiHandler_ServeHTTP(t *testing.T) {
	mux := http.NewServeMux()
	api := newAPIHandler(mux)

	s := httptest.NewServer(api)
	defer s.Close()

	client := req.C().SetTimeout(2 * time.Second).DevMode()

	var resp response
	res, err := client.R().
		SetHeader("Accept", "application/json").
		SetHeader("Accept-Language", "en-US,en;q=0.5").
		SetHeader("X-Real-IP", "172.27.29.34").
		SetHeader("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:50.0) Gecko/20100101 Firefox/50.0").
		SetResult(&resp).
		Get(s.URL + "/api/whoami")
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError() {
		t.Fatal(res.Err)
	}

	expected := response{
		IPAddress: "172.27.29.34",
		Language:  "en-US,en;q=0.5",
		Software:  "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:50.0) Gecko/20100101 Firefox/50.0",
	}
	if expected != resp {
		t.Fatalf("expected %+v, got %+v", expected, resp)
	}
}
