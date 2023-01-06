package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var storedURLs []string

type apiHandler struct{}

type shorturlResp struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

type errorResp struct {
	Error string `json:"error"`
}

func (h *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimPrefix(r.URL.Path, "/api/shorturl")

	if r.Method == http.MethodPost && url == "" {
		if err := r.ParseForm(); err != nil {
			writeErrorResp(w, err)
			return
		}
		original := r.Form.Get("url")
		if original == "" {
			writeErrorResp(w, errors.New("empty url"))
			return
		}

		short, err := shortByOriginalURL(original)
		if err != nil {
			writeErrorResp(w, err)
			return
		}

		writeShorturlResp(w, original, short)
		return
	}

	if r.Method == http.MethodGet {
		url = strings.TrimPrefix(url, "/")
		originalURL, err := originalByShortURL(url)
		if err != nil {
			writeErrorResp(w, err)
			return
		}
		http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
		return
	}
}

func shortByOriginalURL(originalURL string) (string, error) {
	if _, err := http.Get(originalURL); err != nil {
		return "", errors.New("url not exist")
	}

	for i, storedURL := range storedURLs {
		if storedURL == originalURL {
			return strconv.Itoa(i), nil
		}
	}
	storedURLs = append(storedURLs, originalURL)
	return strconv.Itoa(len(storedURLs) - 1), nil
}

func originalByShortURL(shortURL string) (string, error) {
	id, err := strconv.Atoi(shortURL)
	if err != nil {
		return "", errors.New("unknown format")
	}

	if id < 0 || id >= len(storedURLs) {
		return "", errors.New("too large")
	}

	return storedURLs[id], nil
}

func writeShorturlResp(w http.ResponseWriter, originalURL, shortURL string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := &shorturlResp{
		OriginalURL: originalURL,
		ShortURL:    shortURL,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding success json: %v\n", err)
	}
}

func writeErrorResp(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	resp := errorResp{Error: err.Error()}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding error json: %v\n", err)
	}
}
