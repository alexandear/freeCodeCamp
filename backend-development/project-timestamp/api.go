package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var allowedDateFormats = []string{
	"2006-01-02",
	"02 January 2006, GMT",
}

type apiHandler struct {
	clock clock
}

type respSuccess struct {
	UnixMs int64  `json:"unix"`
	UTC    string `json:"utc"`
}

type respError struct {
	Error string `json:"error"`
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dateParam := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/api"), "/")
	if dateParam == "" {
		now := h.clock.Now()
		writeRespSuccess(w, now)
		return
	}

	unixMs, err := strconv.Atoi(dateParam)
	if err == nil {
		writeRespSuccess(w, time.UnixMilli(int64(unixMs)))
		return
	}

	for _, df := range allowedDateFormats {
		date, derr := time.Parse(df, dateParam)
		if derr == nil {
			writeRespSuccess(w, date)
			return
		}
	}

	writeRespError(w)
}

func writeRespSuccess(w http.ResponseWriter, d time.Time) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := respSuccess{
		UnixMs: toUnixMs(d),
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

func toUnixMs(d time.Time) int64 {
	return d.UnixNano() / int64(time.Millisecond)
}
