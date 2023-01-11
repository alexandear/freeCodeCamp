package http

import (
	"context"

	"messageboard/api"
	"messageboard/thread"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 --config=types.cfg.yaml ../api/openapi.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 --config=server.cfg.yaml ../api/openapi.yaml

var _ api.StrictServerInterface = &Server{}

type Server struct {
	threadServ *thread.Service
}

func NewServer(threadServ *thread.Service) *Server {
	return &Server{
		threadServ: threadServ,
	}
}

func (s *Server) CreateThread(ctx context.Context, req api.CreateThreadRequestObject,
) (api.CreateThreadResponseObject, error) {
	res, err := s.threadServ.CreateThread(ctx, thread.CreateThreadParam{
		Board:           req.Board,
		DeletedPassword: req.Body.DeletePassword,
	})
	if err != nil {
		return nil, err
	}

	_ = res
	return api.CreateThread200JSONResponse{}, nil
}
