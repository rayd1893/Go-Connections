package consumer

import "github.com/99minutos/shipments-snapshot-service/pkg/database/mongo"

type Config struct {
	SnapshotsDB mongo.Config `env:", prefix=SS_"`

	Env  string `env:"ENV, default=debug"`
	Port string `env:"PORT, default=8080"`
}
