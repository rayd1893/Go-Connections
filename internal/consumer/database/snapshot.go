package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/99minutos/shipments-snapshot-service/pkg/database/mongo"
	"github.com/99minutos/shipments-snapshot-service/pkg/logging"
)

const (
	shipmentsCollection = "shipments"
	eventsCollection    = "events"
)

type SnapshotsDB struct {
	db *mongo.DB
}

func New(db *mongo.DB) *SnapshotsDB {
	return &SnapshotsDB{
		db: db,
	}
}

func (s *SnapshotsDB) SaveData(ctx context.Context, data map[string]interface{}) error {
	if _, ok := data["trackingId"]; !ok {
		return fmt.Errorf("trackingId is required")
	}
	if id := data["trackingId"].(string); id == "" {
		return fmt.Errorf("trackingId is required")
	}

	logger := logging.FromContext(ctx)
	logger.Debugw("saving event data", "data", data)

	if err := s.db.InTx(ctx, func(sc mongodb.SessionContext) error {
		opts := options.Update()
		opts.SetUpsert(true)

		filter := bson.M{"id": data["trackingId"]}

		events := s.db.DB.Collection(eventsCollection)
		if _, err := events.UpdateOne(sc, filter, bson.M{"$set": data["event"]}, opts); err != nil {
			return err
		}

		shipments := s.db.DB.Collection(shipmentsCollection)
		if _, err := shipments.UpdateOne(sc, filter, bson.M{"$set": data["shipment"]}, opts); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to update data: %w", err)
	}
	return nil
}
