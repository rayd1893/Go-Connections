package databasemanager

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/99minutos/shipments-snapshot-service/pkg/database/mongo"
	"github.com/99minutos/shipments-snapshot-service/pkg/database/postgres"
	"github.com/99minutos/shipments-snapshot-service/pkg/logging"
)

const (
	databaseNotFoundErr          = "%q database with name %q not found"
	databaseAlreadyRegisteredErr = "database %q is already registered"
)

type DBManager struct {
	postgres     map[string]*postgres.DB
	postgresLock sync.RWMutex

	mongo     map[string]*mongo.DB
	mongoLock sync.RWMutex
}

func registerDatabase(ctx context.Context, dbManager *DBManager, cfg interface{}) error {
	switch ctl := cfg.(type) {
	case postgres.Config:
		return dbManager.registerPostgres(ctx, ctl)
	case mongo.Config:
		return dbManager.registerMongo(ctx, ctl)
	}
	return nil
}

func NewFromEnv(ctx context.Context, cfg interface{}) (*DBManager, error) {
	m := &DBManager{
		postgres: make(map[string]*postgres.DB),
		mongo:    make(map[string]*mongo.DB),
	}

	// check if the configuration is a pointer or just an instance of a configuration
	// maybe we may just ensure that any configuration needs to be a pointer?... whatever
	v := reflect.ValueOf(cfg)
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.CanInterface() {
			vf := f.Interface()
			if err := registerDatabase(ctx, m, vf); err != nil {
				return nil, fmt.Errorf("failed to register database: %w", err)
			}
		}
	}

	return m, nil
}

func (m *DBManager) registerPostgres(ctx context.Context, cfg postgres.Config) error {
	logger := logging.FromContext(ctx)

	db, err := postgres.NewFromEnv(ctx, &cfg)
	logger.Infow("connecting to postgres", "config", cfg)
	if err != nil {
		return err
	}

	m.postgresLock.Lock()
	defer m.postgresLock.Unlock()

	connName := cfg.DatabaseConfig().ConnName
	if _, ok := m.postgres[connName]; ok {
		return fmt.Errorf(databaseAlreadyRegisteredErr, connName)
	}

	m.postgres[connName] = db

	return nil
}

func (m *DBManager) GetPostgres(name string) (*postgres.DB, error) {
	m.postgresLock.RLock()
	defer m.postgresLock.RUnlock()

	db, ok := m.postgres[name]
	if !ok {
		return nil, fmt.Errorf(databaseNotFoundErr, "postgres", name)
	}

	return db, nil
}

func (m *DBManager) registerMongo(ctx context.Context, cfg mongo.Config) error {
	logger := logging.FromContext(ctx)

	db, err := mongo.NewFromEnv(ctx, &cfg)
	logger.Infow("connecting to mongo", "config", cfg)
	if err != nil {
		return err
	}

	m.mongoLock.Lock()
	defer m.mongoLock.Unlock()

	connName := cfg.DatabaseConfig().ConnName
	if _, ok := m.postgres[connName]; ok {
		return fmt.Errorf(databaseAlreadyRegisteredErr, connName)
	}

	m.mongo[connName] = db

	return nil
}

func (m *DBManager) GetMongo(name string) (*mongo.DB, error) {
	m.mongoLock.RLock()
	defer m.mongoLock.RUnlock()

	db, ok := m.mongo[name]
	if !ok {
		return nil, fmt.Errorf(databaseNotFoundErr, "mongo", name)
	}

	return db, nil
}

func (m *DBManager) Close(ctx context.Context) []error {
	var errors []error

	for _, db := range m.postgres {
		if err := db.Close(ctx); err != nil {
			errors = append(errors, err)
		}
	}

	for _, db := range m.mongo {
		if err := db.Close(ctx); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
