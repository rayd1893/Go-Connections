package setup

import (
	"context"
	"fmt"

	"github.com/99minutos/shipments-snapshot-service/internal/databasemanager"
	"github.com/99minutos/shipments-snapshot-service/internal/serverenv"
	"github.com/99minutos/shipments-snapshot-service/pkg/logging"
	"github.com/99minutos/shipments-snapshot-service/pkg/secrets"
	"github.com/sethvargo/go-envconfig"
)

type SecretManagerConfigProvider interface {
	SecretManagerConfig() *secrets.Config
}

func Setup(ctx context.Context, config interface{}) (*serverenv.ServerEnv, error) {
	return SetupWith(ctx, config, envconfig.OsLookuper())
}

func SetupWith(ctx context.Context, config interface{}, l envconfig.Lookuper) (*serverenv.ServerEnv, error) {
	logger := logging.FromContext(ctx)

	var mutatorFuncs []envconfig.MutatorFunc

	var serverEnvOpts []serverenv.Option

	var sm secrets.SecretManager
	if provider, ok := config.(SecretManagerConfigProvider); ok {
		logger.Info("configuring secret manager")

		smConfig := provider.SecretManagerConfig()
		if err := envconfig.ProcessWith(ctx, smConfig, l, mutatorFuncs...); err != nil {
			return nil, fmt.Errorf("error processing secret manager config: %w", err)
		}

		var err error
		sm, err = secrets.SecretManagerFor(ctx, smConfig)
		if err != nil {
			return nil, fmt.Errorf("unable to create secret manager: %w", err)
		}

		mutatorFuncs = append(mutatorFuncs, secrets.Resolver(sm, smConfig))

		serverEnvOpts = append(serverEnvOpts, serverenv.WithSecretManager(sm))

		logger.Infow("secret manager configured", "config", smConfig)
	}

	if err := envconfig.ProcessWith(ctx, config, l, mutatorFuncs...); err != nil {
		return nil, fmt.Errorf("error loading environment variables: %w", err)
	}
	logger.Infow("provided", "config", config)

	dbManager, err := databasemanager.NewFromEnv(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error creating database manager: %w", err)
	}
	logger.Infof("configuring database manager")

	serverEnvOpts = append(serverEnvOpts, serverenv.WithDatabaseManager(dbManager))

	return serverenv.New(ctx, serverEnvOpts...), nil
}
