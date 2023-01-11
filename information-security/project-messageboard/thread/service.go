package thread

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type CreateThreadRes struct {
	ThreadID   string
	Text       string
	CreatedOn  time.Time
	BumpedOn   time.Time
	IsReported bool
	Replies    []string
}

func NewService(db *mongo.Database) *Service {
	return &Service{
		threads: db.Collection(ThreadsCollection),
	}
}

func (s *Service) CreateThread(ctx context.Context, param CreateThreadParam) (CreateThreadRes, error) {
	now := time.Now().UTC()
	createdOn := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.UTC)

	res, err := s.threads.InsertOne(ctx, bson.D{
		{"board", param.Board},
		{"text", param.Text},
		{"created_on", createdOn},
		{"bumped_on", createdOn},
		{"is_reported", false},
	})
	if err != nil {
		return CreateThreadRes{}, fmt.Errorf("insert one: %w", err)
	}

	return CreateThreadRes{
		ThreadID:   res.InsertedID.(primitive.ObjectID).Hex(),
		Text:       param.Text,
		CreatedOn:  createdOn,
		BumpedOn:   createdOn,
		IsReported: false,
		Replies:    []string{},
	}, nil
}
