package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		return nil, fmt.Errorf("find: %w", err)
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

func (s *userService) CreateExercise(ctx context.Context, userID, description string, durationMin int, date time.Time,
) (user, exercise, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return user{}, exercise{}, fmt.Errorf("object id from hex: %w", err)
	}

	duration := time.Duration(durationMin) * time.Minute

	if date.IsZero() {
		date = time.Now()
	}

	u, err := s.findUserByID(ctx, objectID)
	if err != nil {
		return user{}, exercise{}, fmt.Errorf("find user by id: %w", err)
	}

	if _, err := s.exerciseColl.InsertOne(ctx, bson.D{
		{"user_id", userID},
		{"description", description},
		{"duration", duration},
		{"date", date},
	}); err != nil {
		return user{}, exercise{}, fmt.Errorf("insert one: %w", err)
	}

	return u, exercise{Description: description, Duration: duration, Date: date}, nil
}

func (s *userService) findUserByID(ctx context.Context, objectID primitive.ObjectID) (user, error) {
	res := s.userColl.FindOne(ctx, bson.M{"_id": objectID})
	if res.Err() != nil {
		return user{}, fmt.Errorf("find one: %w", res.Err())
	}
	var dbUser struct {
		ID       string `bson:"_id"`
		Username string `bson:"username"`
	}
	if err := res.Decode(&dbUser); err != nil {
		return user{}, fmt.Errorf("decode user: %w", err)
	}

	return user{ID: dbUser.ID, Username: dbUser.Username}, nil
}

func (s *userService) Logs(ctx context.Context, userID string, from, to time.Time, limit int) (user, []exercise, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return user{}, nil, fmt.Errorf("object id from hex: %w", err)
	}

	u, err := s.findUserByID(ctx, objectID)
	if err != nil {
		return user{}, nil, fmt.Errorf("find user by id: %w", err)
	}

	opts := options.Find()
	opts.SetLimit(int64(limit))
	filter := bson.D{{"user_id", userID}}
	if !from.IsZero() {
		filter = append(filter, bson.E{"date", bson.D{{"$gte", from}}})
	}
	if !to.IsZero() {
		filter = append(filter, bson.E{"date", bson.D{{"$lte", to}}})
	}

	cursor, err := s.exerciseColl.Find(ctx, filter, opts)
	if err != nil {
		return user{}, nil, fmt.Errorf("find: %w", err)
	}

	var dbExercises []struct {
		UserID      string        `bson:"user_id"`
		Description string        `bson:"description"`
		Duration    time.Duration `bson:"duration"`
		Date        time.Time     `bson:"date"`
	}
	if err := cursor.All(ctx, &dbExercises); err != nil {
		return user{}, nil, fmt.Errorf("all: %w", err)
	}

	exercises := make([]exercise, 0, len(dbExercises))
	for _, dbEx := range dbExercises {
		exercises = append(exercises, exercise{
			Description: dbEx.Description,
			Duration:    dbEx.Duration,
			Date:        dbEx.Date,
		})
	}

	return u, exercises, nil
}
