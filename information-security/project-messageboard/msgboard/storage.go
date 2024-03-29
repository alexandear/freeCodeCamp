package msgboard

import (
	"time"
)

const (
	ThreadsCollection = "threads"
	RepliesCollection = "replies"
)

type storageThread struct {
	ThreadID       string    `bson:"_id"`
	Text           string    `bson:"text"`
	CreatedOn      time.Time `bson:"created_on"`
	BumpedOn       time.Time `bson:"bumped_on"`
	IsReported     bool      `bson:"is_reported"`
	DeletePassword []byte    `bson:"delete_password"`
	ReplyCount     int       `bson:"reply_count"`
}

func (t *storageThread) ToThread(storageReplies []storageReply) ThreadRes {
	replies := make([]ReplyRes, 0, len(storageReplies))
	for _, storageReply := range storageReplies {
		replies = append(replies, storageReply.ToReply())
	}

	return ThreadRes{
		ThreadID:   t.ThreadID,
		Text:       t.Text,
		CreatedOn:  t.CreatedOn,
		BumpedOn:   t.BumpedOn,
		IsReported: t.IsReported,
		ReplyCount: t.ReplyCount,
		Replies:    replies,
	}
}

type storageReply struct {
	ReplyID        string    `bson:"_id"`
	ThreadID       string    `bson:"thread_id"`
	Board          string    `bson:"board"`
	Text           string    `bson:"text"`
	CreatedOn      time.Time `bson:"created_on"`
	DeletePassword []byte    `bson:"delete_password"`
	IsReported     bool      `bson:"is_reported"`
}

func (r *storageReply) ToReply() ReplyRes {
	return ReplyRes{
		ReplyID:        r.ReplyID,
		ThreadID:       r.ThreadID,
		Text:           r.Text,
		CreatedOn:      r.CreatedOn,
		DeletePassword: r.DeletePassword,
	}
}
