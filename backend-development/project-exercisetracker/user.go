package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type user struct {
	ID       string
	Username string
}

type exercise struct {
	Description string
	Duration    time.Duration
	Date        time.Time
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

func (s *userService) CreateUser(ctx context.Context, username string) (string, error) {
	res, err := s.userColl.InsertOne(ctx, bson.D{{"username", username}})
	if err != nil {
		return "", fmt.Errorf("insert one: %w", err)
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *userService) AllUsers(ctx context.Context) ([]user, error) {
	cursor, err := s.userColl.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, fmt.Errorf("find one: %w", err)
	}

	var dbUsers []struct {
		ID       string `bson:"_id"`
		Username string `bson:"username"`
	}
	if err := cursor.All(ctx, &dbUsers); err != nil {
		return nil, fmt.Errorf("all: %w", err)
	}

	users := make([]user, 0, len(dbUsers))
	for _, dbUser := range dbUsers {
		users = append(users, user{
			ID:       dbUser.ID,
			Username: dbUser.Username,
		})
	}
	return users, nil
}

func (s *userService) CreateExercise(ctx context.Context, userID, description, durationStr, dateStr string,
) (user, exercise, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return user{}, exercise{}, fmt.Errorf("object id from hex: %w", err)
	}

	durationMin, err := strconv.Atoi(durationStr)
	if err != nil {
		return user{}, exercise{}, fmt.Errorf("duration invalid: %w", err)
	}
	duration := time.Duration(durationMin) * time.Minute

	var date time.Time
	if dateStr != "" {
		date = time.Now().UTC()
		d, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return user{}, exercise{}, fmt.Errorf("date parse: %w", err)
		}
		date = d
	} else {
		date = time.Now().UTC()
	}

	res := s.userColl.FindOne(ctx, bson.M{"_id": objectID})
	if res.Err() != nil {
		return user{}, exercise{}, fmt.Errorf("find one: %w", res.Err())
	}
	var dbUser struct {
		Username string `bson:"username"`
	}
	if err := res.Decode(&dbUser); err != nil {
		return user{}, exercise{}, fmt.Errorf("decode user: %w", err)
	}

	if _, err := s.exerciseColl.InsertOne(ctx, bson.D{
		{"user_id", userID},
		{"description", description},
		{"duration", duration},
		{"date", date},
	}); err != nil {
		return user{}, exercise{}, fmt.Errorf("insert one: %w", err)
	}

	return user{
			ID:       userID,
			Username: dbUser.Username,
		},
		exercise{
			Description: description,
			Duration:    duration,
			Date:        date,
		}, nil
}
