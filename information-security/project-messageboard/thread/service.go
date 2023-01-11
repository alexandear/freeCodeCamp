package thread

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	threads *mongo.Collection
}

type CreateThreadParam struct {
	Board           string
	DeletedPassword string
}

type CreateThreadRes struct {
	ThreadID string
}

func NewService(db *mongo.Database) *Service {
	return &Service{
		threads: db.Collection("threads"),
	}
}

func (s *Service) CreateThread(ctx context.Context, param CreateThreadParam) (CreateThreadRes, error) {
	return CreateThreadRes{}, nil
}
