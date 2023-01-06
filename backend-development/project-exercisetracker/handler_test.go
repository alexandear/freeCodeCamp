package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

var client = &http.Client{
	Timeout: 200 * time.Millisecond,
}

func TestHandlerServeStaticContent(t *testing.T) {
	h := newHandler(echo.New())

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
