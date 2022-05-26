package serverenv

import (
	"context"

	"github.com/99minutos/shipments-snapshot-service/internal/databasemanager"
	"github.com/99minutos/shipments-snapshot-service/pkg/secrets"
)

type ServerEnv struct {
	databaseManager *databasemanager.DBManager
	secretManager   secrets.SecretManager
}

type Option func(*ServerEnv) *ServerEnv

func New(ctx context.Context, opts ...Option) *ServerEnv {
	env := &ServerEnv{}

	for _, f := range opts {
		env = f(env)
	}

	return env
}

func WithDatabaseManager(m *databasemanager.DBManager) Option {
	return func(s *ServerEnv) *ServerEnv {
		s.databaseManager = m
		return s
	}
}

func WithSecretManager(sm secrets.SecretManager) Option {
	return func(s *ServerEnv) *ServerEnv {
		s.secretManager = sm
		return s
	}
}

func (s *ServerEnv) DatabaseManager() *databasemanager.DBManager {
	return s.databaseManager
}

func (s *ServerEnv) Close(ctx context.Context) error {
	if s == nil {
		return nil
	}

	if s.databaseManager != nil {
		s.databaseManager.Close(ctx)
	}

	return nil
}
