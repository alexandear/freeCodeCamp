package http

import (
	"net/http"

	"messageboard/api"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 --config=types.cfg.yaml ../api/openapi.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 --config=server.cfg.yaml ../api/openapi.yaml

var _ api.ServerInterface = &Server{}

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) CreateThread(w http.ResponseWriter, r *http.Request) {
}
