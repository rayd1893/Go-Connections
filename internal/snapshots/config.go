package snapshots

import (
	"time"

	"github.com/99minutos/shipments-snapshot-service/pkg/database/mongo"
)

type Config struct {
	FSMDatabase mongo.Config `env:", prefix=FSM_"`

	Env     string        `env:"ENV, default=debug"`
	Port    string        `env:"PORT, default=8080"`
	Timeout time.Duration `env:"RPC_TIMEOUT, default=10m"`
}
