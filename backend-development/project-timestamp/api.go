package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "2006-01-02"

type apiHandler struct {
	clock clock
}

type respSuccess struct {
	UnixMs *int64 `json:"unix,omitempty"`
	UTC    string `json:"utc,omitempty"`
}

type respError struct {
	Error string `json:"error"`
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dateParam := strings.TrimPrefix(r.URL.Path, "/api")
	if dateParam == "" || dateParam == "/" {
		now := h.clock.Now()
		ms := unixMs(now)
		writeRespSuccess(w, &ms, now)
		return
	}

	unixMs, err := strconv.ParseInt(dateParam, 10, 64)
	if err == nil {
		writeRespSuccess(w, &unixMs, time.UnixMilli(unixMs))
		return
	}

	date, err := time.Parse(dateFormat, dateParam)
	if err != nil {
		writeRespError(w)
		return
	}

	writeRespSuccess(w, nil, date)
}

func writeRespSuccess(w http.ResponseWriter, unixMs *int64, d time.Time) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := respSuccess{
		UnixMs: unixMs,
		UTC:    d.Format(http.TimeFormat),
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding success json: %v\n", err)
	}
}

func writeRespError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	resp := respError{Error: "Invalid Date"}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding error json: %v\n", err)
	}
}

func unixMs(d time.Time) int64 {
	return d.UnixNano() / int64(time.Millisecond)
}
