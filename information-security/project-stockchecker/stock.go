package main

import (
	"context"
	"crypto/sha256"
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
	StockTwo   string
	IfLike     bool
	RemoteAddr string
}

type storageStock struct {
	LikesCount int `bson:"likes_count"`
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

	ipHash, err := hashIP(param.RemoteAddr)
	if err != nil {
		return StockData{}, err
	}

	q, err := quote.Get(stock)
	if err != nil {
		return StockData{}, fmt.Errorf("quote get: %w", err)
	}
	price := q.Ask

	var likesInc int
	upsert := options.Update().SetUpsert(true)
	if param.IfLike {
		stockPerIP := param.Stock + "-" + ipHash
		update := bson.D{{"$set", bson.D{{"_id", stockPerIP}}}}
		res, err := s.stockPerIPs.UpdateByID(ctx, stockPerIP, update, upsert)
		if err != nil {
			return StockData{}, fmt.Errorf("update stock per ips: %w", err)
		}
		if res.MatchedCount == 0 {
			likesInc = 1
		}
	}

	update := bson.D{{"$inc", bson.D{{"likes_count", likesInc}}}}
	if _, err = s.stocks.UpdateByID(ctx, stock, update, upsert); err != nil {
		return StockData{}, fmt.Errorf("update stock: %w", err)
	}

	var ss storageStock
	if err := s.stocks.FindOne(ctx, bson.M{"_id": stock}).Decode(&ss); err != nil {
		return StockData{}, fmt.Errorf("find stock: %w", err)
	}

	return StockData{
		Price:      price,
		LikesCount: ss.LikesCount,
	}, nil
}

func (s *StockService) StockDatas(ctx context.Context, stocks []string) ([]StockData, error) {
	stockDatas := make([]StockData, len(stocks))

	find := make(bson.D, len(stocks))
	for i, stock := range stocks {
		q, err := quote.Get(stock)
		if err != nil {
			return nil, fmt.Errorf("quote for %s get: %w", stock, err)
		}
		stockDatas[i].Price = q.Ask

		find = append(find, bson.E{"_id", stock})
	}

	cursor, err := s.stocks.Find(ctx, find)
	if err != nil {
		return nil, fmt.Errorf("find stock: %w", err)
	}

	var dbStocks []storageStock
	if err := cursor.All(ctx, &dbStocks); err != nil {
		return nil, err
	}

	return stockDatas, nil
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
