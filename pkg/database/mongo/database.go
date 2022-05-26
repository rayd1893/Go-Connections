package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func (db *DB) InTx(ctx context.Context, f func(sc mongo.SessionContext) error) error {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := db.Conn.StartSession()
	if err != nil {
		return fmt.Errorf("starting session: %w", err)
	}
	defer session.EndSession(ctx)

	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		err := sc.StartTransaction(txOpts)
		if err != nil {
			return fmt.Errorf("starting transaction: %w", err)
		}

		if err := f(sc); err != nil {
			if err1 := session.AbortTransaction(ctx); err1 != nil {
				return fmt.Errorf("rolling back transaction: %v (original error: %w)", err1, err)
			}
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if err := session.CommitTransaction(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}
