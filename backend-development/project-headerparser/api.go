package main

import (
	"encoding/json"
	"net/http"

	"github.com/rs/cors"
)

type apiHandler struct {
}

func newAPIHandler(mux *http.ServeMux) http.Handler {
	mux.Handle("/api/whoami", &apiHandler{})
	return cors.AllowAll().Handler(mux)
}

type response struct {
	IPAddress string `json:"ipaddress"`
	Language  string `json:"language"`
	Software  string `json:"software"`
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response{
		IPAddress: r.RemoteAddr,
		Language:  r.Header.Get("Accept-Language"),
		Software:  r.Header.Get("User-Agent"),
	})
}
