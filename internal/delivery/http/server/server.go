package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vagonaizer/effective-mobile/subscription-service/internal/config"
	"github.com/vagonaizer/effective-mobile/subscription-service/pkg/logger"
)

type Server struct {
	config                 config.ServerConfig
	router                 *gin.Engine
	httpServer             *http.Server
	logger                 *logger.Logger
	readTimeout            time.Duration
	writeTimeout           time.Duration
	idleTimeout            time.Duration
	shutdownTimeout        time.Duration
	enableGracefulShutdown bool
	healthCheck            func(ctx context.Context) error
}

func New(opts ...Option) *Server {
	server := &Server{
		readTimeout:            30 * time.Second,
		writeTimeout:           30 * time.Second,
		idleTimeout:            60 * time.Second,
		shutdownTimeout:        30 * time.Second,
		enableGracefulShutdown: true,
	}

	for _, opt := range opts {
		opt(server)
	}

	if server.logger == nil {
		defaultLogger, _ := logger.NewDefaultLogger()
		server.logger = defaultLogger
	}

	server.setupHTTPServer()
	return server
}

func (s *Server) setupHTTPServer() {
	s.httpServer = &http.Server{
		Addr:           s.config.Address(),
		Handler:        s.router,
		ReadTimeout:    s.readTimeout,
		WriteTimeout:   s.writeTimeout,
		IdleTimeout:    s.idleTimeout,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
}

func (s *Server) Start() error {
	s.logger.Info("starting http server",
		zap.String("address", s.config.Address()),
		zap.Duration("read_timeout", s.readTimeout),
		zap.Duration("write_timeout", s.writeTimeout))

	if s.healthCheck != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.healthCheck(ctx); err != nil {
			s.logger.Error("health check failed", zap.Error(err))
			return err
		}
		s.logger.Info("health check passed")
	}

	if s.enableGracefulShutdown {
		return s.startWithGracefulShutdown()
	}

	s.logger.Info("server started successfully", zap.String("address", s.config.Address()))
	return s.httpServer.ListenAndServe()
}

func (s *Server) startWithGracefulShutdown() error {
	go func() {
		s.logger.Info("server started successfully", zap.String("address", s.config.Address()))
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal("server startup failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	s.logger.Info("shutdown signal received", zap.String("signal", sig.String()))

	return s.Shutdown()
}

func (s *Server) Shutdown() error {
	s.logger.Info("shutting down server gracefully",
		zap.Duration("timeout", s.shutdownTimeout))

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("server forced to shutdown", zap.Error(err))
		return err
	}

	s.logger.Info("server shutdown completed")
	return nil
}

func (s *Server) GetHTTPServer() *http.Server {
	return s.httpServer
}

func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

func (s *Server) GetConfig() config.ServerConfig {
	return s.config
}

func (s *Server) SetupTimeouts() {
	if s.config.ReadTimeout > 0 {
		s.readTimeout = time.Duration(s.config.ReadTimeout) * time.Second
	}
	if s.config.WriteTimeout > 0 {
		s.writeTimeout = time.Duration(s.config.WriteTimeout) * time.Second
	}
	if s.config.IdleTimeout > 0 {
		s.idleTimeout = time.Duration(s.config.IdleTimeout) * time.Second
	}
	s.setupHTTPServer()
}
