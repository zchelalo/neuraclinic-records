package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zchelalo/neuraclinic-records/pkg/bootstrap"
	"go.uber.org/zap"
)

func main() {
	cfg, err := bootstrap.LoadConfig(".env")
	if err != nil {
		panic(err)
	}

	logger := bootstrap.GetLogger()
	defer bootstrap.SyncLogger()

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := bootstrap.InitApp(rootCtx, logger, cfg)
	if err != nil {
		logger.Fatal("cannot initialize application", zap.Error(err))
	}

	errCh := make(chan error, 2)
	go func() {
		logger.Info("grpc server starting", zap.Int("port", cfg.Port))
		errCh <- app.Server.Start()
	}()
	go func() {
		logger.Info("rabbitmq consumer starting", zap.String("queue", cfg.RabbitMQQueue))
		errCh <- app.Consumer.Run(rootCtx)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-sigs:
		logger.Info("signal received, shutting down", zap.String("signal", sig.String()))
	case err := <-errCh:
		if err != nil {
			logger.Error("worker stopped", zap.Error(err))
		}
	}

	cancel()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = app.Cleanup(ctx)
	logger.Info("shutdown complete")
}
