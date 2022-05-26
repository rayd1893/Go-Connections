package postgres

import "github.com/99minutos/shipments-snapshot-service/pkg/database"

type Config struct {
	Base database.Config
}

func (c *Config) DatabaseConfig() *database.Config {
	return &c.Base
}
