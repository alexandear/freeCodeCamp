package msgboard

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Transaction struct {
	ctx     context.Context
	session mongo.Session
}

func NewTransaction(ctx context.Context, dbClient *mongo.Client) (*Transaction, error) {
	session, err := dbClient.StartSession()
	if err != nil {
		return nil, fmt.Errorf("start session: %w", err)
	}
	return &Transaction{session: session, ctx: ctx}, nil
}

func (t *Transaction) Close() {
	t.session.EndSession(t.ctx)
}

func (t *Transaction) Start(fn func(ctx mongo.SessionContext) (any, error),
) (any, error) {
	txnOptions := options.Transaction().SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
	return t.session.WithTransaction(t.ctx, fn, txnOptions)
}
