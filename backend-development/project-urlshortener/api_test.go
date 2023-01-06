package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestApiHandlerPost(t *testing.T) {
	api := &apiHandler{}

	s := httptest.NewServer(api)
	defer s.Close()

	res, err := http.PostForm(s.URL+"/api/shorturl", url.Values{"url": {"https://google.com"}})
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	var resp shorturlResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	expected := shorturlResp{
		OriginalURL: "https://google.com",
		ShortURL:    "0",
	}
	if expected != resp {
		t.Fatalf("expected %+v, got %+v", expected, resp)
	}
}

func TestApiHandlerPostError(t *testing.T) {
	api := &apiHandler{}

	s := httptest.NewServer(api)
	defer s.Close()

	res, err := http.PostForm(s.URL+"/api/shorturl", url.Values{"url": {"ftp:/john-doe.invalidTLD"}})
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	var resp errorResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	expected := errorResp{
		Error: "invalid url",
	}
	if expected != resp {
		t.Fatalf("expected %+v, got %+v", expected, resp)
	}
}

func TestApiHandlerGet(t *testing.T) {
	api := &apiHandler{}

	s := httptest.NewServer(api)
	defer s.Close()

	_, err := http.PostForm(s.URL+"/api/shorturl", url.Values{"url": {"https://google.com"}})
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.Get(s.URL + "/api/shorturl/0")
	if err != nil {
		t.Fatal(err)
	}

	if http.StatusOK != res.StatusCode {
		t.Fatal("status must be OK")
	}

	if http.StatusMovedPermanently != res.Request.Response.StatusCode {
		t.Fatal("status must be Moved Permanently")
	}

	if "https://google.com" != res.Request.Header.Get("Referer") {
		t.Fatal("wrong redirect url")
	}
}
