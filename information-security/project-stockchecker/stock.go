package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"

	"github.com/piquette/finance-go/quote"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StockData struct {
	Price      float64
	LikesCount int
}

type StockDataParam struct {
	Stock      string
	IfLike     bool
	RemoteAddr string
}

type storageStock struct {
	Stock      string  `bson:"stock"`
	Price      float64 `bson:"price"`
	LikesCount int     `bson:"likes_count"`
}

type StockService struct {
	stocks      *mongo.Collection
	stockPerIPs *mongo.Collection
}

func NewStockService(db *mongo.Database) *StockService {
	return &StockService{
		stocks:      db.Collection("stocks"),
		stockPerIPs: db.Collection("stock_per_ips"),
	}
}

func (s *StockService) StockData(ctx context.Context, param StockDataParam) (StockData, error) {
	stock := param.Stock
	var incLike bool

	ipHash, err := hashIP(param.RemoteAddr)
	if err != nil {
		return StockData{}, err
	}

	if param.IfLike {
		stockPerIP := param.Stock + "-" + ipHash
		update := bson.D{{"$set", bson.D{{"_id", stockPerIP}}}}
		res, err := s.stockPerIPs.UpdateByID(ctx, stockPerIP, update, options.Update().SetUpsert(true))
		if err != nil {
			return StockData{}, fmt.Errorf("update by id: %w", err)
		}
		if res.MatchedCount == 0 {
			incLike = true
		}
	}

	var ss storageStock
	ferr := s.stocks.FindOne(ctx, bson.M{"stock": stock}).Decode(&ss)
	if errors.Is(ferr, mongo.ErrNoDocuments) {
		q, err := quote.Get(stock)
		if err != nil {
			return StockData{}, fmt.Errorf("quote get: %w", err)
		}

		var likesCount int
		if incLike {
			likesCount = 1
		}

		ss := &storageStock{
			Stock:      stock,
			Price:      q.Ask,
			LikesCount: likesCount,
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
	if ferr != nil {
		return StockData{}, fmt.Errorf("find one: %w", ferr)
	}

	return StockData{
		Price:      ss.Price,
		LikesCount: ss.LikesCount,
	}, nil
}

func hashIP(remoteAddr string) (string, error) {
	ipPort := strings.Split(remoteAddr, ":")
	if len(ipPort) != 2 {
		return "", fmt.Errorf("remote addr must have format 'ip:port'")
	}

	sha := sha256.New()
	sha.Write([]byte(ipPort[0]))
	return string(sha.Sum(nil)), nil
}
