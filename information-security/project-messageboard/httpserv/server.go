package httpserv

import (
	"context"
	"fmt"

	"messageboard/api"
	"messageboard/msgboard"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 --config=types.cfg.yaml ../api/openapi.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 --config=server.cfg.yaml ../api/openapi.yaml

var _ api.StrictServerInterface = &Server{}

type Server struct {
	msgServ *msgboard.Service
}

func NewServer(msgServ *msgboard.Service) *Server {
	return &Server{
		msgServ: msgServ,
	}
}

func (s *Server) GetThreads(ctx context.Context, req api.GetThreadsRequestObject) (api.GetThreadsResponseObject, error) {
	threads, err := s.msgServ.Threads(ctx, req.Board)
	if err != nil {
		return nil, fmt.Errorf("get msgboard: %w", err)
	}

	res := make(api.GetThreads200JSONResponse, 0, len(threads))
	for _, th := range threads {
		resThread := api.Thread{
			Id:        th.ThreadID,
			BumpedOn:  th.BumpedOn,
			CreatedOn: th.CreatedOn,
			Text:      th.Text,
			Replies:   []api.Reply{},
		}

		for _, r := range th.Replies {
			resThread.Replies = append(resThread.Replies, api.Reply{
				Id:   r.ReplyID,
				Text: r.Text,
			})
		}

		res = append(res, resThread)
	}

	return res, nil
}

func (s *Server) CreateThread(ctx context.Context, req api.CreateThreadRequestObject,
) (api.CreateThreadResponseObject, error) {
	body := req.JSONBody
	if req.FormdataBody != nil {
		body = req.FormdataBody
	}
	err := s.msgServ.CreateThread(ctx, msgboard.CreateThreadParam{
		Board:          req.Board,
		Text:           body.Text,
		DeletePassword: body.DeletePassword,
	})
	if err != nil {
		return nil, fmt.Errorf("create msgboard: %w", err)
	}

	return api.CreateThread200Response{}, nil
}

func (s *Server) CreateReply(ctx context.Context, req api.CreateReplyRequestObject) (api.CreateReplyResponseObject, error) {
	body := req.JSONBody
	if req.FormdataBody != nil {
		body = req.FormdataBody
	}
	err := s.msgServ.CreateReply(ctx, msgboard.CreateReplyParam{
		ThreadID:       body.ThreadId,
		Board:          req.Board,
		Text:           body.Text,
		DeletePassword: body.DeletePassword,
	})
	if err != nil {
		return nil, fmt.Errorf("create reply: %w", err)
	}

	return api.CreateReply200Response{}, nil
}
