package httpserv

import (
	"context"
	"fmt"

	"messageboard/api"
	"messageboard/msgboard"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 --config=types.cfg.yaml ../api/openapi.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 --config=server.cfg.yaml ../api/openapi.yaml

const (
	DeleteThreadRespSuccess           = "success"
	DeleteThreadRespIncorrectPassword = "incorrect password"
)

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
	for _, thread := range threads {
		res = append(res, toAPIThread(thread))
	}

	return res, nil
}

func (s *Server) CreateThread(ctx context.Context, req api.CreateThreadRequestObject,
) (api.CreateThreadResponseObject, error) {
	body := req.JSONBody
	if req.FormdataBody != nil {
		body = req.FormdataBody
	}
	threadID, err := s.msgServ.CreateThread(ctx, msgboard.CreateThreadParam{
		Board:          req.Board,
		Text:           body.Text,
		DeletePassword: body.DeletePassword,
	})
	if err != nil {
		return nil, fmt.Errorf("create msgboard: %w", err)
	}

	return api.CreateThread200TextResponse(threadID), nil
}

func (s *Server) DeleteThread(ctx context.Context, req api.DeleteThreadRequestObject) (api.DeleteThreadResponseObject, error) {
	ifDeleted, err := s.msgServ.DeleteThread(ctx, msgboard.DeleteThreadParam{
		Board:          req.Board,
		ThreadID:       req.Body.ThreadId,
		DeletePassword: req.Body.DeletePassword,
	})
	if err != nil {
		return nil, fmt.Errorf("delete thread: %w", err)
	}

	if ifDeleted {
		return api.DeleteThread200TextResponse(DeleteThreadRespSuccess), nil
	}
	return api.DeleteThread200TextResponse(DeleteThreadRespIncorrectPassword), nil
}

func (s *Server) CreateReply(ctx context.Context, req api.CreateReplyRequestObject) (api.CreateReplyResponseObject, error) {
	body := req.JSONBody
	if req.FormdataBody != nil {
		body = req.FormdataBody
	}
	replyID, err := s.msgServ.CreateReply(ctx, msgboard.CreateReplyParam{
		ThreadID:       body.ThreadId,
		Board:          req.Board,
		Text:           body.Text,
		DeletePassword: body.DeletePassword,
	})
	if err != nil {
		return nil, fmt.Errorf("create reply: %w", err)
	}

	return api.CreateReply200TextResponse(replyID), nil
}

func (s *Server) GetReplies(ctx context.Context, req api.GetRepliesRequestObject) (api.GetRepliesResponseObject, error) {
	threadID := req.Params.ThreadId
	thread, err := s.msgServ.Thread(ctx, req.Board, threadID)
	if err != nil {
		return nil, fmt.Errorf("thread: %w", err)
	}

	return api.GetReplies200JSONResponse(toAPIThread(thread)), nil
}

func toAPIThread(thread msgboard.ThreadRes) api.Thread {
	replies := make([]api.Reply, 0, len(thread.Replies))
	for _, r := range thread.Replies {
		replies = append(replies, api.Reply{
			Id:        r.ReplyID,
			Text:      r.Text,
			CreatedOn: r.CreatedOn,
		})
	}
	return api.Thread{
		Id:        thread.ThreadID,
		BumpedOn:  thread.BumpedOn,
		CreatedOn: thread.CreatedOn,
		Replies:   replies,
		Text:      thread.Text,
	}
}
