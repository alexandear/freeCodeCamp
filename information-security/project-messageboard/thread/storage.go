package thread

import (
	"time"
)

type storageThread struct {
	ThreadID   string    `bson:"_id"`
	Text       string    `bson:"text"`
	CreatedOn  time.Time `bson:"created_on"`
	BumpedOn   time.Time `bson:"bumped_on"`
	IsReported bool      `bson:"is_reported"`
}

func (t *storageThread) ToThread(replies []ReplyRes) ThreadRes {
	return ThreadRes{
		ThreadID:   t.ThreadID,
		Text:       t.Text,
		CreatedOn:  t.CreatedOn,
		BumpedOn:   t.BumpedOn,
		IsReported: t.IsReported,
		Replies:    replies,
	}
}

type storageReply struct {
	ReplyID  string `bson:"_id"`
	ThreadID string `bson:"thread_id"`
	Text     string `bson:"text"`
}

func (r *storageReply) ToReply() ReplyRes {
	return ReplyRes{
		ReplyID:  r.ReplyID,
		ThreadID: r.ThreadID,
		Text:     r.Text,
	}
}
