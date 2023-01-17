package msgboard

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type transaction struct {
	ctx    context.Context
	client *mongo.Client
}

func NewTransaction(ctx context.Context, client *mongo.Client) *transaction {
	return &transaction{
		ctx:    ctx,
		client: client,
	}
}

func (t *transaction) Start(fn func(ctx mongo.SessionContext) (any, error)) error {
	session, err := t.client.StartSession()
	if err != nil {
		return fmt.Errorf("start session: %w", err)
	}
	defer session.EndSession(t.ctx)

	txnOptions := options.Transaction().SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
	if _, err := session.WithTransaction(t.ctx, fn, txnOptions); err != nil {
		return fmt.Errorf("execute transaction: %w", err)
	}

	return nil
}
