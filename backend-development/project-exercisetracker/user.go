package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type user struct {
	ID       string `bson:"_id" json:"_id"`
	Username string `bson:"username" json:"username"`
}

type exercise struct {
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Date        string `json:"date"`
	User        user
}

type userService struct {
	userColl     *mongo.Collection
	exerciseColl *mongo.Collection
}

func newUserService(db *mongo.Database) *userService {
	return &userService{
		userColl:     db.Collection("users"),
		exerciseColl: db.Collection("exercises"),
	}
}

func (s *userService) CreateUser(ctx context.Context, username string) (user, error) {
	res, err := s.userColl.InsertOne(ctx, bson.D{{"username", username}})
	if err != nil {
		return user{}, fmt.Errorf("insert one: %w", err)
	}

	return user{
		ID:       res.InsertedID.(primitive.ObjectID).Hex(),
		Username: username,
	}, nil
}

func (s *userService) AllUsers(ctx context.Context) ([]user, error) {
	cursor, err := s.userColl.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, fmt.Errorf("find one: %w", err)
	}

	var users []user
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("all: %w", err)
	}

	return users, nil
}

func (s *userService) CreateExercise(
	ctx context.Context, userID, description string, duration int, date time.Time,
) (exercise, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return exercise{}, fmt.Errorf("object id from hex: %w", err)
	}

	res := s.userColl.FindOne(ctx, bson.M{"_id": objectID})
	if res.Err() != nil {
		return exercise{}, fmt.Errorf("find one: %w", res.Err())
	}
	var u user
	if err := res.Decode(&u); err != nil {
		return exercise{}, fmt.Errorf("decode user: %w", err)
	}

	if _, err := s.exerciseColl.InsertOne(ctx, bson.D{
		{"user_id", userID},
		{"description", description},
		{"duration", duration},
		{"date", date},
	}); err != nil {
		return exercise{}, fmt.Errorf("insert one: %w", err)
	}

	return exercise{
		Description: description,
		Duration:    duration,
		Date:        date.Format("Mon Jan 01 2006"),
		User:        u,
	}, nil
}
