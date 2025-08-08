package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/vagonaizer/effective-mobile/subscription-service/internal/config"
	"github.com/vagonaizer/effective-mobile/subscription-service/pkg/logger"
)

type App struct {
	deps   *Dependencies
	logger *logger.Logger
}

func New(configPath string) (*App, error) {
	cfg := config.NewConfig()
	if err := cfg.Load(configPath); err != nil {
		return nil, err
	}

	loggerConfig := logger.Config{
		Level:       cfg.Logger.Level,
		Development: cfg.Logger.Development,
		Encoding:    cfg.Logger.Encoding,
	}

	log, err := logger.NewLogger(loggerConfig)
	if err != nil {
		return nil, err
	}

	log.Info("application starting",
		zap.String("version", "1.0.0"),
		zap.String("environment", getEnvironment(cfg.Logger.Development)))

	deps, err := NewDependencies(*cfg, log)
	if err != nil {
		log.Error("failed to initialize dependencies", zap.Error(err))
		return nil, err
	}

	return &App{
		deps:   deps,
		logger: log,
	}, nil
}

func (a *App) Run() error {
	a.logger.Info("starting subscription service",
		zap.String("address", a.deps.Config.Server.Address()))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		if err := a.deps.Server.Start(); err != nil {
			errChan <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		a.logger.Error("server error", zap.Error(err))
		return err
	case sig := <-quit:
		a.logger.Info("shutdown signal received", zap.String("signal", sig.String()))
		return a.shutdown(ctx)
	}
}

func (a *App) shutdown(ctx context.Context) error {
	a.logger.Info("gracefully shutting down application")

	if err := a.deps.Server.Shutdown(); err != nil {
		a.logger.Error("server shutdown error", zap.Error(err))
		return err
	}

	if err := a.deps.Close(); err != nil {
		a.logger.Error("dependencies cleanup error", zap.Error(err))
		return err
	}

	a.logger.Info("application shutdown completed successfully")
	return nil
}

func (a *App) GetDependencies() *Dependencies {
	return a.deps
}

func (a *App) GetLogger() *logger.Logger {
	return a.logger
}

func (a *App) HealthCheck(ctx context.Context) error {
	return a.deps.Database.HealthCheck(ctx)
}

func getEnvironment(development bool) string {
	if development {
		return "development"
	}
	return "production"
}
