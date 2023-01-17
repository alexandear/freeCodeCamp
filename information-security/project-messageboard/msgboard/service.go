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
	"golang.org/x/crypto/bcrypt"
)

const (
	maxReturnedThreadsCount = 10
	maxReturnedRepliesCount = 3
)

type Service struct {
	dbClient *mongo.Client
	threads  *mongo.Collection
	replies  *mongo.Collection
}

func NewService(db *mongo.Database) *Service {
	return &Service{
		dbClient: db.Client(),
		threads:  db.Collection(ThreadsCollection),
		replies:  db.Collection(RepliesCollection),
	}
}

func (s *Service) Threads(ctx context.Context, board string) ([]ThreadRes, error) {
	opts := options.Find().SetLimit(maxReturnedThreadsCount).SetSort(bson.M{"bumped_on": -1})
	cursor, err := s.threads.Find(ctx, bson.M{"board": board}, opts)
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
	deletePassword, err := hashPassword(param.DeletePassword)
	if err != nil {
		return "", fmt.Errorf("make hash: %w", err)
	}

	createdOn := now()
	threadID := primitive.NewObjectID()

	if _, err := s.threads.InsertOne(ctx, bson.D{
		{"_id", threadID},
		{"board", param.Board},
		{"text", param.Text},
		{"created_on", createdOn},
		{"bumped_on", createdOn},
		{"delete_password", deletePassword},
		{"is_reported", false},
		{"reply_count", 0},
	}); err != nil {
		return "", fmt.Errorf("insert one: %w", err)
	}

	return threadID.Hex(), nil
}

func (s *Service) DeleteThread(ctx context.Context, param DeleteThreadParam) (bool, error) {
	threadObjectID, err := primitive.ObjectIDFromHex(param.ThreadID)
	if err != nil {
		return false, fmt.Errorf("wrong object id: %w", err)
	}

	res, err := NewTransaction(ctx, s.dbClient).Start(func(ctx mongo.SessionContext) (any, error) {
		var dbThread storageThread
		err := s.threads.FindOne(ctx, bson.M{"_id": threadObjectID}).Decode(&dbThread)
		if err != nil {
			return false, fmt.Errorf("find one: %w", err)
		}

		if !compareHash(dbThread.DeletePassword, param.DeletePassword) {
			return false, nil
		}

		if _, err = s.threads.DeleteOne(ctx, bson.M{"_id": threadObjectID}); err != nil {
			return false, fmt.Errorf("delete one: %w", err)
		}

		if _, err := s.replies.DeleteMany(ctx, bson.M{"thread_id": param.ThreadID}); err != nil {
			return false, fmt.Errorf("delete replies: %w", err)
		}

		return true, nil
	})
	if err != nil {
		return false, err
	}

	return res.(bool), nil
}

func (s *Service) ReportThread(ctx context.Context, board, threadID string) error {
	threadObjectID, err := primitive.ObjectIDFromHex(threadID)
	if err != nil {
		return fmt.Errorf("wrong object id: %w", err)
	}

	res, err := s.threads.UpdateOne(ctx,
		bson.D{{"_id", threadObjectID}, {"board", board}},
		bson.M{"$set": bson.M{"is_reported": true}})
	if err != nil {
		return fmt.Errorf("update one: %w", err)
	}

	if res.MatchedCount != 1 {
		return fmt.Errorf("thread not found")
	}

	return nil
}

func (s *Service) CreateReply(ctx context.Context, param CreateReplyParam) (string, error) {
	threadObjectID, err := primitive.ObjectIDFromHex(param.ThreadID)
	if err != nil {
		return "", fmt.Errorf("wrong object id: %w", err)
	}

	deletePassword, err := hashPassword(param.DeletePassword)
	if err != nil {
		return "", fmt.Errorf("make hash: %w", err)
	}

	n := now()
	update := bson.D{
		{"$set", bson.M{"bumped_on": n}},
		{"$inc", bson.M{"reply_count": 1}},
	}

	res, err := NewTransaction(ctx, s.dbClient).Start(func(ctx mongo.SessionContext) (any, error) {
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
			{"board", param.Board},
			{"text", param.Text},
			{"created_on", n},
			{"delete_password", deletePassword},
			{"is_reported", false},
		})
		if err != nil {
			return "", fmt.Errorf("insert one: %w", err)
		}

		return replyID.Hex(), nil
	})
	if err != nil {
		return "", err
	}

	return res.(string), nil
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

func (s *Service) DeleteReply(ctx context.Context, param DeleteReplyParam) (bool, error) {
	replyObjectID, err := primitive.ObjectIDFromHex(param.ReplyID)
	if err != nil {
		return false, fmt.Errorf("wrong object id: %w", err)
	}

	res, err := NewTransaction(ctx, s.dbClient).Start(func(ctx mongo.SessionContext) (any, error) {
		filter := bson.D{
			{"_id", replyObjectID},
			{"thread_id", param.ThreadID},
			{"board", param.Board},
		}

		var dbReply storageReply
		err = s.replies.FindOne(ctx, filter).Decode(&dbReply)
		if err != nil {
			return false, fmt.Errorf("find one: %w", err)
		}

		if !compareHash(dbReply.DeletePassword, param.DeletePassword) {
			return false, nil
		}

		if _, err := s.replies.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"text": "[deleted]"}}); err != nil {
			return false, fmt.Errorf("delete one: %w", err)
		}

		return true, nil
	})
	if err != nil {
		return false, err
	}

	return res.(bool), nil
}

func (s *Service) ReportReply(ctx context.Context, board, threadID, replyID string) error {
	replyObjectID, err := primitive.ObjectIDFromHex(replyID)
	if err != nil {
		return fmt.Errorf("wrong object id: %w", err)
	}

	res, err := s.replies.UpdateOne(ctx, bson.D{
		{"_id", replyObjectID},
		{"thread_id", threadID},
		{"board", board},
	}, bson.D{{"$set", bson.M{"is_reported": true}}})
	if err != nil {
		return fmt.Errorf("update one: %w", err)
	}

	if res.ModifiedCount != 1 {
		return fmt.Errorf("reply not found")
	}

	return nil
}

func now() time.Time {
	return time.Now().UTC()
}

func hashPassword(password string) ([]byte, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return pw, nil
}

func compareHash(hashedPassword []byte, password string) bool {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)) == nil
}
