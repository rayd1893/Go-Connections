package snapshots

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/99minutos/shipments-snapshot-service/internal/serverenv"
	"github.com/99minutos/shipments-snapshot-service/internal/snapshots/database"
	"github.com/99minutos/shipments-snapshot-service/internal/snapshots/model"
	"github.com/99minutos/shipments-snapshot-service/pkg/logging"
	"github.com/99minutos/shipments-snapshot-service/pkg/pb"
)

const (
	missingDBErr = "missing %q database in database manager"
)

var errorFailedToRetrieveShipment = errors.New("failed to retrieve shipment")

var _ pb.ShipmentsSnapshotServer = (*Server)(nil)

func NewServer(env *serverenv.ServerEnv, config *Config) (pb.ShipmentsSnapshotServer, error) {
	if env.DatabaseManager() == nil {
		return nil, fmt.Errorf("missing database manager in server environment")
	}

	mongoconn, err := env.DatabaseManager().GetMongo(StateMachineDatabase)
	if err != nil {
		return nil, fmt.Errorf(missingDBErr, StateMachineDatabase)
	}

	return &Server{
		db:     database.New(mongoconn),
		env:    env,
		config: config,
	}, nil
}

type Server struct {
	pb.UnimplementedShipmentsSnapshotServer

	db     *database.StateMachineDB
	env    *serverenv.ServerEnv
	config *Config
}

func (s Server) GetShipment(ctx context.Context, req *pb.GetShipmentRequest) (*pb.Shipment, error) {
	logger := logging.FromContext(ctx).Named("snapshots.GetShipment")

	ctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	shipment, err := s.db.FindOne(ctx, req.GetTrackingId())
	if err != nil {
		logger.Errorw("failed to retrieve shipment", "error", err)
		return nil, errorFailedToRetrieveShipment
	}

	res, err := model.TransformShipmentToProto(shipment)
	if err != nil {
		logger.Errorw("failed to transform to proto", "error", err)
		return nil, errorFailedToRetrieveShipment
	}
	return res, nil
}

func buildCriteria(req *pb.GetShipmentsRequest) (*database.IterateShipmentsCriteria, error) {
	criteria := &database.IterateShipmentsCriteria{}

	start := req.GetIntervalDates().GetStart()
	if start != "" {
		t, err := time.Parse(time.RFC3339, start)
		if err != nil {
			return nil, fmt.Errorf("invalid start date: %w", err)
		}
		criteria.FromStart = t
	}

	end := req.GetIntervalDates().GetFinish()
	if end != "" {
		t, err := time.Parse(time.RFC3339, end)
		if err != nil {
			return nil, fmt.Errorf("invalid end date: %w", err)
		}
		criteria.ToEnd = t
	}

	criteria.FromClient = req.GetClientId()
	criteria.OnlyThisDeliveryType = req.GetDeliveryType()
	criteria.FromStatus = req.GetStatus()
	return criteria, nil
}

func (s Server) GetShipments(ctx context.Context, req *pb.GetShipmentsRequest) (*pb.GetShipmentsResponse, error) {
	logger := logging.FromContext(ctx).Named("snapshots.GetShipments")

	ctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	criteria, err := buildCriteria(req)
	if err != nil {
		logger.Error("failed to build criteria", "error", err)
		return nil, errors.New("invalid request")
	}

	var shipments []*pb.Shipment

	opt := options.Find()
	ps := int64(req.GetPageCursor().GetPageSize())
	opt.SetLimit(ps)
	opt.SetSkip(int64(req.GetPageCursor().GetPage()) * ps)
	if err := s.db.IterateShipments(ctx, *criteria, func(shipment *model.Shipment) error {
		s, err := model.TransformShipmentToProto(shipment)
		if err != nil {
			return fmt.Errorf("failed to transform shipment to proto: %w", err)
		}
		shipments = append(shipments, s)
		return nil
	}, opt); err != nil {
		logger.Errorw("failed to iterate shipments", "error", err)
		return nil, errors.New("failed to retrieve shipments")
	}

	return &pb.GetShipmentsResponse{Shipments: shipments}, nil
}

func (s Server) GetEvents(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	logger := logging.FromContext(ctx).Named("snapshots.GetEvents")

	ctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	var events []*pb.Event
	if err := s.db.IterateEvents(ctx, req.GetTrackingId(), func(event *model.Event) error {
		s, err := model.TransformEventToProto(event)
		if err != nil {
			return fmt.Errorf("failed to transform event to proto: %w", err)
		}
		events = append(events, s)
		return nil
	}); err != nil {
		logger.Errorw("failed to iterate events", "error", err)
		return nil, errors.New("failed to retrieve events")
	}

	return &pb.GetEventsResponse{Events: events}, nil
}
