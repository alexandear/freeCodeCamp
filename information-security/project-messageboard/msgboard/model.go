package msgboard

import (
	"time"
)

type CreateThreadParam struct {
	Board          string
	Text           string
	DeletePassword string
}

type DeleteThreadParam struct {
	Board          string
	ThreadID       string
	DeletePassword string
}

type CreateReplyParam struct {
	Board          string
	Text           string
	DeletePassword string
	ThreadID       string
}

type DeleteReplyParam struct {
	Board          string
	ThreadID       string
	ReplyID        string
	DeletePassword string
}

type ThreadRes struct {
	ThreadID   string
	Text       string
	CreatedOn  time.Time
	BumpedOn   time.Time
	IsReported bool
	Replies    []ReplyRes
	ReplyCount int
}

type ReplyRes struct {
	ReplyID        string
	ThreadID       string
	Text           string
	CreatedOn      time.Time
	DeletePassword []byte
}
