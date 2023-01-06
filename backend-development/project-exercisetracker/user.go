package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type user struct {
	ID       string `bson:"_id" json:"_id"`
	Username string `bson:"username" json:"username"`
}

type userService struct {
	coll *mongo.Collection
}

func newUserService(db *mongo.Database) *userService {
	return &userService{coll: db.Collection("users")}
}

func (s *userService) CreateUser(ctx context.Context, username string) (user, error) {
	res, err := s.coll.InsertOne(ctx, bson.D{{"username", username}})
	if err != nil {
		return user{}, fmt.Errorf("insert one: %w", err)
	}

	return user{
		ID:       res.InsertedID.(primitive.ObjectID).Hex(),
		Username: username,
	}, nil
}

func (s *userService) AllUsers(ctx context.Context) ([]user, error) {
	cursor, err := s.coll.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, fmt.Errorf("find one: %w", err)
	}

	var users []user
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("all: %w", err)
	}

	return users, nil
}
