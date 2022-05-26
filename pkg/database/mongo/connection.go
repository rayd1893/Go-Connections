package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/99minutos/shipments-snapshot-service/pkg/database"
)

var _ database.Database = (*DB)(nil)

type DB struct {
	Conn *mongo.Client
	DB   *mongo.Database
}

func NewFromEnv(ctx context.Context, cfg *Config) (*DB, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.DatabaseConfig().ConnectionURL()))
	if err != nil {
		return nil, fmt.Errorf("failed to create client connection: %w", err)
	}

	db := &DB{Conn: client}
	err = db.Connect(ctx)
	db.DB = db.Conn.Database(cfg.DatabaseConfig().Name)
	return db, err
}

func (db *DB) Connect(ctx context.Context) error {
	if err := db.Conn.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to mongo database: %w", err)
	}
	return nil
}

func (db *DB) Close(ctx context.Context) error {
	if err := db.Conn.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to close mongo database: %w", err)
	}
	return nil
}
