package consumer

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/99minutos/shipments-snapshot-service/internal/consumer/database"
	"github.com/99minutos/shipments-snapshot-service/internal/jsonutil"
	"github.com/99minutos/shipments-snapshot-service/internal/serverenv"
	consumerapi "github.com/99minutos/shipments-snapshot-service/pkg/api/v1"
	"github.com/99minutos/shipments-snapshot-service/pkg/logging"
	"github.com/99minutos/shipments-snapshot-service/pkg/middlewares"
	"github.com/99minutos/shipments-snapshot-service/pkg/render"
)

const (
	missingDBErr = "missing %q database in database manager"
)

type Server struct {
	config *Config
	db     *database.SnapshotsDB
	env    *serverenv.ServerEnv
	h      *render.Renderer
}

func NewServer(config *Config, env *serverenv.ServerEnv) (*Server, error) {
	if env.DatabaseManager() == nil {
		return nil, fmt.Errorf("missing database manager in server environment")
	}

	mongoconn, err := env.DatabaseManager().GetMongo(ShipmentSnapshotDatabase)
	if err != nil {
		return nil, fmt.Errorf(missingDBErr, ShipmentSnapshotDatabase)
	}

	return &Server{
		db:     database.New(mongoconn),
		config: config,
		env:    env,
		h:      render.NewRenderer(),
	}, nil
}

func (s *Server) Routes() *mux.Router {
	r := mux.NewRouter()
	r.Use(middlewares.Recovery())
	r.Use(middlewares.PopulateRequestID())

	r.Handle("/", s.handleEvent()).Methods(http.MethodPost)

	return r
}

func (s *Server) handleEvent() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger := logging.FromContext(ctx).Named("consumer.handleEvent")

		var event consumerapi.Event
		code, err := jsonutil.Unmarshal(w, r, &event)
		if err != nil {
			logger.Errorw("failed to unmarshal event", "error", err)
			s.h.RenderJSON(w, code, nil)
			return
		}

		data, err := event.GetDecodedData()
		if err != nil {
			logger.Errorw("failed to decode event data", "error", err)
			s.h.RenderJSON(w, http.StatusInternalServerError, nil)
			return
		}

		if err := s.db.SaveData(ctx, data); err != nil {
			logger.Errorw("failed to save shipment", "error", err)
			s.h.RenderJSON(w, http.StatusInternalServerError, nil)
			return
		}

		s.h.RenderJSON(w, code, nil)
	})
}
