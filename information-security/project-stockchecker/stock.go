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
	"golang.org/x/sync/errgroup"
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

type StockDataParam2 struct {
	Stock1     string
	Stock2     string
	IfLike     bool
	RemoteAddr string
}

type storageStock struct {
	Stock      string `bson:"_id"`
	LikesCount int    `bson:"likes_count"`
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

func (s *StockService) StockDataAndLike(ctx context.Context, param StockDataParam) (StockData, error) {
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

func (s *StockService) StockDataAndLike2(ctx context.Context, param StockDataParam2) ([2]StockData, error) {
	g, gctx := errgroup.WithContext(ctx)

	var (
		price1 float64
		likes1 int
	)
	g.Go(func() error {
		p := StockDataParam{
			Stock:      param.Stock1,
			IfLike:     param.IfLike,
			RemoteAddr: param.RemoteAddr,
		}
		sd, err := s.StockDataAndLike(gctx, p)
		if err != nil {
			return fmt.Errorf("stock data for %s: %w", param.Stock1, err)
		}
		price1 = sd.Price
		likes1 = sd.LikesCount
		return nil
	})
	var (
		price2 float64
		likes2 int
	)
	g.Go(func() error {
		p := StockDataParam{
			Stock:      param.Stock2,
			IfLike:     param.IfLike,
			RemoteAddr: param.RemoteAddr,
		}
		sd, err := s.StockDataAndLike(gctx, p)
		if err != nil {
			return fmt.Errorf("stock data for %s: %w", param.Stock2, err)
		}
		price2 = sd.Price
		likes2 = sd.LikesCount
		return nil
	})
	if err := g.Wait(); err != nil {
		return [2]StockData{}, err
	}

	return [2]StockData{
		{
			Price:      price1,
			LikesCount: likes1,
		},
		{
			Price:      price2,
			LikesCount: likes2,
		},
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
