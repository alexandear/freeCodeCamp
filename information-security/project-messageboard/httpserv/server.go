package httpserv

import (
	"context"
	"fmt"

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
	body := req.JSONBody
	if req.FormdataBody != nil {
		body = req.FormdataBody
	}
	res, err := s.threadServ.CreateThread(ctx, thread.CreateThreadParam{
		Board:          req.Board,
		Text:           body.Text,
		DeletePassword: body.DeletePassword,
	})
	if err != nil {
		return nil, fmt.Errorf("create thread: %w", err)
	}

	return api.CreateThread200JSONResponse{
		Id:        res.ThreadID,
		BumpedOn:  res.BumpedOn,
		CreatedOn: res.CreatedOn,
		Replies:   res.Replies,
		Reported:  res.IsReported,
		Text:      res.Text,
	}, nil
}
