package msgboard

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	maxReturnedThreadsCount = 10
	maxReturnedRepliesCount = 3
)

type Service struct {
	threads *mongo.Collection
	replies *mongo.Collection
}

type CreateThreadParam struct {
	Board          string
	Text           string
	DeletePassword string
}

type CreateReplyParam struct {
	Board          string
	Text           string
	DeletePassword string
	ThreadID       string
}

type ThreadRes struct {
	ThreadID   string
	Text       string
	CreatedOn  time.Time
	BumpedOn   time.Time
	IsReported bool
	Replies    []ReplyRes
}

type ReplyRes struct {
	ReplyID  string
	ThreadID string
	Text     string
}

func NewService(db *mongo.Database) *Service {
	return &Service{
		threads: db.Collection(ThreadsCollection),
		replies: db.Collection(RepliesCollection),
	}
}

func (s *Service) Threads(ctx context.Context, board string) ([]ThreadRes, error) {
	opts := options.Find().SetLimit(maxReturnedThreadsCount).SetSort(bson.D{{"bumped_on", -1}})
	cursor, err := s.threads.Find(ctx, bson.D{{"board", board}}, opts)
	if err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}

	dbThreads := make([]storageThread, 0, maxReturnedThreadsCount)
	if err := cursor.All(ctx, &dbThreads); err != nil {
		return nil, fmt.Errorf("cursor all: %w", err)
	}

	threads := make([]ThreadRes, 0, len(dbThreads))
	for _, dbThread := range dbThreads {
		replies, err := s.RepliesForThread(ctx, dbThread.ThreadID, maxReturnedRepliesCount)
		if err != nil {
			return nil, fmt.Errorf("replies for msgboard=%s: %w", dbThread.ThreadID, err)
		}

		threads = append(threads, dbThread.ToThread(replies))
	}

	return threads, nil
}

func (s *Service) Thread(ctx context.Context, board, threadID string) (ThreadRes, error) {
	threadObjID, err := primitive.ObjectIDFromHex(threadID)
	if err != nil {
		return ThreadRes{}, fmt.Errorf("object id from hex: %w", err)
	}

	var dbThread storageThread
	err = s.threads.FindOne(ctx, bson.D{{"board", board}, {"_id", threadObjID}}).Decode(&dbThread)
	if err != nil {
		return ThreadRes{}, fmt.Errorf("find one thread: %w", err)
	}

	replies, err := s.RepliesForThread(ctx, threadID, 0)
	if err != nil {
		return ThreadRes{}, fmt.Errorf("find replies: %w", err)
	}

	return dbThread.ToThread(replies), nil
}

func (s *Service) CreateThread(ctx context.Context, param CreateThreadParam) (string, error) {
	createdOn := now()

	threadID := primitive.NewObjectID()
	_, err := s.threads.InsertOne(ctx, bson.D{
		{"_id", threadID},
		{"board", param.Board},
		{"text", param.Text},
		{"created_on", createdOn},
		{"bumped_on", createdOn},
		{"is_reported", false},
	})
	if err != nil {
		return "", fmt.Errorf("insert one: %w", err)
	}

	return threadID.Hex(), nil
}

func (s *Service) CreateReply(ctx context.Context, param CreateReplyParam) (string, error) {
	threadObjectID, err := primitive.ObjectIDFromHex(param.ThreadID)
	if err != nil {
		return "", fmt.Errorf("wrong object id: %w", err)
	}

	n := now()
	update := bson.D{{"$set", bson.D{{"bumped_on", n}}}}
	_, err = s.threads.UpdateByID(ctx, threadObjectID, update)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", fmt.Errorf("board not found: %w", err)
	}
	if err != nil {
		return "", fmt.Errorf("find one and update err: %w", err)
	}

	replyID := primitive.NewObjectID()

	_, err = s.replies.InsertOne(ctx, bson.D{
		{"_id", replyID},
		{"thread_id", param.ThreadID},
		{"text", param.Text},
		{"created_on", n},
		{"delete_password", param.DeletePassword},
	})
	if err != nil {
		return "", fmt.Errorf("insert one: %w", err)
	}

	return replyID.Hex(), nil
}

func (s *Service) RepliesForThread(ctx context.Context, threadID string, limit int) ([]ReplyRes, error) {
	cursor, err := s.replies.Find(ctx, bson.D{{"thread_id", threadID}}, options.Find().SetLimit(int64(limit)))
	if err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}

	dbReplies := make([]storageReply, 0, limit)
	if err := cursor.All(ctx, &dbReplies); err != nil {
		return nil, fmt.Errorf("cursor all: %w", err)
	}

	replies := make([]ReplyRes, 0, len(dbReplies))
	for _, dbReply := range dbReplies {
		replies = append(replies, dbReply.ToReply())
	}

	return replies, nil
}

func now() time.Time {
	return time.Now().UTC()
}
