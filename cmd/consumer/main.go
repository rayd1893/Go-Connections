package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/99minutos/shipments-snapshot-service/internal/consumer"
	"github.com/99minutos/shipments-snapshot-service/internal/setup"
	"github.com/99minutos/shipments-snapshot-service/pkg/logging"
	"github.com/99minutos/shipments-snapshot-service/pkg/server"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	logger := logging.NewLoggerFromEnv()
	ctx = logging.WithLogger(ctx, logger)

	defer func() {
		done()
		if r := recover(); r != nil {
			logger.Fatalw("application panic", "panic", r)
		}
	}()

	err := realMain(ctx)
	done()

	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("successful shutdown")
}

func realMain(ctx context.Context) error {
	logger := logging.FromContext(ctx)

	var config consumer.Config
	env, err := setup.Setup(ctx, &config)
	if err != nil {
		return fmt.Errorf("setup.Setup: %w", err)
	}
	defer env.Close(ctx)

	consumerServer, err := consumer.NewServer(&config, env)
	if err != nil {
		return fmt.Errorf("consumer.NewServer: %w", err)
	}

	srv, err := server.New(config.Port)
	if err != nil {
		return fmt.Errorf("server.New: %w", err)
	}
	logger.Infof("listen on: %s", config.Port)

	return srv.ServeHTTPHandler(ctx, consumerServer.Routes())
}
