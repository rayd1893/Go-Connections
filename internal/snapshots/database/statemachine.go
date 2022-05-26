package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/99minutos/shipments-snapshot-service/internal/snapshots/model"
	"github.com/99minutos/shipments-snapshot-service/pkg/database/mongo"
	"github.com/99minutos/shipments-snapshot-service/pkg/logging"
)

const (
	shipmentsCollection = "shipments"
	eventsCollection    = "events"
)

type StateMachineDB struct {
	db *mongo.DB
}

func New(db *mongo.DB) *StateMachineDB {
	return &StateMachineDB{
		db: db,
	}
}

type IterateShipmentsCriteria struct {
	FromStart            time.Time
	ToEnd                time.Time
	FromClient           string
	OnlyThisDeliveryType string
	FromStatus           string
}

type (
	ShipmentsIteratorFunction func(*model.Shipment) error
	EventsIteratorFunction    func(*model.Event) error
)

func (sm *StateMachineDB) FindOne(ctx context.Context, id string) (*model.Shipment, error) {
	if id == "" {
		return nil, errors.New("id is empty")
	}

	logger := logging.FromContext(ctx).Named("FindOne")
	logger.Debugf("finding a shipment with id: %s", id)

	var shipment model.Shipment
	if err := sm.db.InTx(ctx, func(sc mongodb.SessionContext) error {
		shipments := sm.db.DB.Collection(shipmentsCollection)

		if err := shipments.FindOne(sc, bson.M{"trackingId": id}).Decode(&shipment); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to retrieve shipment: %w", err)
	}

	return &shipment, nil
}

func generateShipmentsFilter(criteria IterateShipmentsCriteria) bson.M {
	filter := make(bson.M)

	if criteria.FromStatus != "" {
		filter["status"] = criteria.FromStatus
	}
	if criteria.FromClient != "" {
		filter["client_id"] = criteria.FromClient
	}
	if criteria.OnlyThisDeliveryType != "" {
		filter["deliveryType"] = criteria.OnlyThisDeliveryType
	}
	if !criteria.ToEnd.IsZero() && !criteria.FromStart.IsZero() {
		filter["createdAt"] = bson.M{
			"$gt": criteria.FromStart,
			"$lt": criteria.ToEnd,
		}
	}
	return filter
}

func (sm *StateMachineDB) IterateShipments(ctx context.Context, criteria IterateShipmentsCriteria, f ShipmentsIteratorFunction, opts ...*options.FindOptions) error {
	filter := generateShipmentsFilter(criteria)

	logger := logging.FromContext(ctx).Named("IterateShipments")
	logger.Debugw("iterator query", "query", filter)

	if err := sm.db.InTx(ctx, func(sc mongodb.SessionContext) error {
		shipments := sm.db.DB.Collection(shipmentsCollection)
		cursor, err := shipments.Find(sc, filter, opts...)
		if err != nil {
			return err
		}
		defer cursor.Close(sc)

		var results []*model.Shipment
		if err := cursor.All(sc, &results); err != nil {
			return err
		}

		for _, shipment := range results {
			if err := f(shipment); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to iterate shipments: %w", err)
	}

	return nil
}

func (sm *StateMachineDB) IterateEvents(ctx context.Context, id string, f EventsIteratorFunction, opts ...*options.FindOptions) error {
	logger := logging.FromContext(ctx).Named("IterateEvents")
	logger.Debugw("iterating events", "tracking_id", id)

	shipment, err := sm.FindOne(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find shipment: %w", err)
	}
	shipmentID, err := primitive.ObjectIDFromHex(shipment.ID)
	if err != nil {
		return fmt.Errorf("failed to convert shipment id to object id: %w", err)
	}

	if err := sm.db.InTx(ctx, func(sc mongodb.SessionContext) error {
		events := sm.db.DB.Collection(eventsCollection)
		cursor, err := events.Find(sc, bson.M{"shipmentId": shipmentID}, opts...)
		if err != nil {
			return err
		}
		defer cursor.Close(sc)

		var results []*model.Event
		if err := cursor.All(sc, &results); err != nil {
			return err
		}

		for _, event := range results {
			if err := f(event); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to iterate events: %w", err)
	}

	return nil
}
