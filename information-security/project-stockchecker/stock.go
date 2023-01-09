package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/piquette/finance-go/quote"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StockData struct {
	Price      float64
	LikesCount int
}

type storageStock struct {
	Stock      string  `bson:"stock"`
	Price      float64 `bson:"price"`
	LikesCount int     `bson:"likes_count"`
}

type StockService struct {
	stocks *mongo.Collection
}

func NewStockService(db *mongo.Database) *StockService {
	return &StockService{
		stocks: db.Collection("stocks"),
	}
}

func (s *StockService) StockData(ctx context.Context, stock string) (StockData, error) {
	res := s.stocks.FindOne(ctx, bson.M{"stock": stock})
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		q, err := quote.Get(stock)
		if err != nil {
			return StockData{}, fmt.Errorf("quote get: %w", err)
		}

		ss := &storageStock{
			Stock:      stock,
			Price:      q.Ask,
			LikesCount: 0,
		}
		_, insErr := s.stocks.InsertOne(ctx, ss)
		if insErr != nil {
			return StockData{}, fmt.Errorf("insert one: %w", err)
		}

		return StockData{
			Price:      ss.Price,
			LikesCount: ss.LikesCount,
		}, nil
	}
	if res.Err() != nil {
		return StockData{}, fmt.Errorf("find one: %w", res.Err())
	}

	var ss storageStock
	if decErr := res.Decode(&ss); decErr != nil {
		return StockData{}, fmt.Errorf("decode: %w", decErr)
	}

	return StockData{
		Price:      ss.Price,
		LikesCount: ss.LikesCount,
	}, nil
}
