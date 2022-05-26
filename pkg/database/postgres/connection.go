package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/99minutos/shipments-snapshot-service/pkg/database"
	"github.com/99minutos/shipments-snapshot-service/pkg/logging"
)

var _ database.Database = (*DB)(nil)

type DB struct {
	Conn *pgxpool.Pool
}

func NewFromEnv(ctx context.Context, cfg *Config) (*DB, error) {
	pgxConfig, err := pgxpool.ParseConfig(dbDSN(cfg))
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	pgxConfig.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		return conn.Ping(ctx) == nil
	}

	pool, err := pgxpool.ConnectConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return &DB{Conn: pool}, nil
}

func (db *DB) Connect(ctx context.Context) error {
	return nil
}

func (db *DB) Close(ctx context.Context) error {
	logger := logging.FromContext(ctx)
	logger.Infof("closing postgres connection")
	db.Conn.Close()
	return nil
}

func dbDSN(cfg *Config) string {
	vals := dbValues(cfg)
	p := make([]string, 0, len(vals))
	for k, v := range vals {
		p = append(p, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(p, " ")
}

func setIfNotEmpty(m map[string]string, key, val string) {
	if val != "" {
		m[key] = val
	}
}

func dbValues(cfg *Config) map[string]string {
	base := cfg.DatabaseConfig()
	p := map[string]string{}
	setIfNotEmpty(p, "dbname", base.Name)
	setIfNotEmpty(p, "user", base.User)
	setIfNotEmpty(p, "host", base.Host)
	setIfNotEmpty(p, "port", base.Port)
	setIfNotEmpty(p, "password", base.Password)
	return p
}
