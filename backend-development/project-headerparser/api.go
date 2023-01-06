package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type apiHandler struct {
}

type respSuccess struct {
	IPAddress string `json:"ipaddress"`
	Language  string `json:"language"`
	Software  string `json:"software"`
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := respSuccess{
		IPAddress: r.RemoteAddr,
		Language:  r.Header.Get("Accept-Language"),
		Software:  r.Header.Get("User-Agent"),
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding success json: %v\n", err)
	}
}
