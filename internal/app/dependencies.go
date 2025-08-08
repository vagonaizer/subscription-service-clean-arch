package app

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/vagonaizer/effective-mobile/subscription-service/internal/config"
	"github.com/vagonaizer/effective-mobile/subscription-service/internal/delivery/http/handlers"
	"github.com/vagonaizer/effective-mobile/subscription-service/internal/delivery/http/middleware"
	"github.com/vagonaizer/effective-mobile/subscription-service/internal/delivery/http/router"
	"github.com/vagonaizer/effective-mobile/subscription-service/internal/delivery/http/server"
	"github.com/vagonaizer/effective-mobile/subscription-service/internal/domain/ports/repository"
	"github.com/vagonaizer/effective-mobile/subscription-service/internal/domain/ports/service"
	"github.com/vagonaizer/effective-mobile/subscription-service/internal/infrastructure/database/postgres"
	infraRepo "github.com/vagonaizer/effective-mobile/subscription-service/internal/infrastructure/database/postgres/repository"
	appService "github.com/vagonaizer/effective-mobile/subscription-service/internal/service"
	"github.com/vagonaizer/effective-mobile/subscription-service/pkg/logger"
)

type Dependencies struct {
	Config config.Config
	Logger *logger.Logger

	Database *postgres.DB

	SubscriptionRepo    repository.SubscriptionRepository
	SubscriptionService service.SubscriptionService

	SubscriptionHandler *handlers.SubscriptionHandler
	HealthHandler       *handlers.HealthHandler

	Router *router.Router
	Server *server.Server
}

func NewDependencies(cfg config.Config, log *logger.Logger) (*Dependencies, error) {
	deps := &Dependencies{
		Config: cfg,
		Logger: log,
	}

	if err := deps.initDatabase(); err != nil {
		return nil, err
	}

	if err := deps.initRepositories(); err != nil {
		return nil, err
	}

	if err := deps.initServices(); err != nil {
		return nil, err
	}

	if err := deps.initHandlers(); err != nil {
		return nil, err
	}

	if err := deps.initRouter(); err != nil {
		return nil, err
	}

	if err := deps.initServer(); err != nil {
		return nil, err
	}

	return deps, nil
}

func (d *Dependencies) initDatabase() error {
	d.Logger.Info("initializing database connection")

	db, err := postgres.New(d.Config.Database, d.Logger)
	if err != nil {
		return err
	}

	d.Database = db
	d.Logger.Info("database connection initialized successfully")
	return nil
}

func (d *Dependencies) initRepositories() error {
	d.Logger.Info("initializing repositories")

	d.SubscriptionRepo = infraRepo.NewSubscriptionRepository(d.Database, d.Logger)

	d.Logger.Info("repositories initialized successfully")
	return nil
}

func (d *Dependencies) initServices() error {
	d.Logger.Info("initializing services")

	d.SubscriptionService = appService.NewSubscriptionService(d.SubscriptionRepo, d.Logger)

	d.Logger.Info("services initialized successfully")
	return nil
}

func (d *Dependencies) initHandlers() error {
	d.Logger.Info("initializing handlers")

	d.SubscriptionHandler = handlers.NewSubscriptionHandler(d.SubscriptionService, d.Logger)

	d.HealthHandler = handlers.NewHealthHandler(d.Logger, func(ctx context.Context) error {
		return d.Database.HealthCheck(ctx)
	})

	d.Logger.Info("handlers initialized successfully")
	return nil
}

func (d *Dependencies) initRouter() error {
	d.Logger.Info("initializing router")

	routerConfig := router.RouterConfig{
		Debug:  d.Config.Logger.Development,
		Logger: d.Logger,
	}

	r := router.New(routerConfig)

	middlewares := []gin.HandlerFunc{
		middleware.CORS(),
		middleware.StructuredLogger(d.Logger),
		middleware.Recovery(d.Logger),
		middleware.ErrorHandler(d.Logger),
	}
	r.SetupMiddleware(middlewares...)

	r.RegisterHealthRoutes()
	r.RegisterAPIRoutes(
		d.SubscriptionHandler,
		d.HealthHandler,
	)
	r.RegisterSwaggerRoutes()

	d.Router = r
	d.Logger.Info("router initialized successfully")
	return nil
}

func (d *Dependencies) initServer() error {
	d.Logger.Info("initializing server")

	d.Server = server.New(
		server.WithConfig(d.Config.Server),
		server.WithLogger(d.Logger),
		server.WithRouter(d.Router.Engine()),
		server.WithGracefulShutdown(),
		server.WithHealthCheck(func(ctx context.Context) error {
			return d.Database.HealthCheck(ctx)
		}),
	)

	d.Server.SetupTimeouts()

	d.Logger.Info("server initialized successfully")
	return nil
}

func (d *Dependencies) Close() error {
	d.Logger.Info("closing dependencies")

	if d.Database != nil {
		d.Database.Close()
	}

	if d.Logger != nil {
		d.Logger.Sync()
	}

	d.Logger.Info("dependencies closed successfully")
	return nil
}
