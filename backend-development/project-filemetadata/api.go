package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

func makeAPIHandler(mux *http.ServeMux) http.Handler {
	mux.Handle("/api/fileanalyse", http.HandlerFunc(uploadFile))
	return cors.AllowAll().Handler(mux)
}

type response struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int64  `json:"size"`
}

func uploadFile(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Maximum upload of 10 MB files
	_ = req.ParseMultipartForm(10 << 20)

	// Get handler for filename, size and headers
	file, handler, err := req.FormFile("upfile")
	if err != nil {
		http.Error(w, "retrieving the file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	defer file.Close()

	_ = json.NewEncoder(w).Encode(&response{
		Name: handler.Filename,
		Type: fmt.Sprintf("%v", handler.Header.Get("Content-Type")),
		Size: handler.Size,
	})
}
