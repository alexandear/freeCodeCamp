package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestApiHandler_ServeHTTP(t *testing.T) {
	api := &apiHandler{}

	s := httptest.NewServer(api)
	defer s.Close()

	req, err := http.NewRequest(http.MethodGet, s.URL+"/api/whoami", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:50.0) Gecko/20100101 Firefox/50.0")
	client := &http.Client{Timeout: time.Second}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	var resp respSuccess
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	expected := respSuccess{
		IPAddress: resp.IPAddress,
		Language:  "en-US,en;q=0.5",
		Software:  "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:50.0) Gecko/20100101 Firefox/50.0",
	}
	if expected != resp {
		t.Fatalf("expected %+v, got %+v", expected, resp)
	}
}
