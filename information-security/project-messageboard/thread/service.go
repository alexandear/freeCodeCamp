package thread

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ThreadsCollection = "threads"
)

type Service struct {
	threads *mongo.Collection
}

type CreateThreadParam struct {
	Board          string
	Text           string
	DeletePassword string
}

type ThreadRes struct {
	ThreadID   string
	Text       string
	CreatedOn  time.Time
	BumpedOn   time.Time
	IsReported bool
	Replies    []string
}

type storageThread struct {
	ThreadID   string    `bson:"_id"`
	Text       string    `bson:"text"`
	CreatedOn  time.Time `bson:"created_on"`
	BumpedOn   time.Time `bson:"bumped_on"`
	IsReported bool      `bson:"is_reported"`
}

func NewService(db *mongo.Database) *Service {
	return &Service{
		threads: db.Collection(ThreadsCollection),
	}
}

func (s *Service) Thread(ctx context.Context, board string) (ThreadRes, error) {
	var thread storageThread
	err := s.threads.FindOne(ctx, bson.D{{"board", board}}).Decode(&thread)
	if err != nil {
		return ThreadRes{}, err
	}

	return ThreadRes{
		ThreadID:   thread.ThreadID,
		Text:       thread.Text,
		CreatedOn:  thread.CreatedOn,
		BumpedOn:   thread.BumpedOn,
		IsReported: thread.IsReported,
		Replies:    []string{},
	}, nil
}

func (s *Service) CreateThread(ctx context.Context, param CreateThreadParam) error {
	now := time.Now().UTC()
	createdOn := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.UTC)

	_, err := s.threads.InsertOne(ctx, bson.D{
		{"board", param.Board},
		{"text", param.Text},
		{"created_on", createdOn},
		{"bumped_on", createdOn},
		{"is_reported", false},
	})
	if err != nil {
		return fmt.Errorf("insert one: %w", err)
	}

	return nil
}
