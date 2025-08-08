package server

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vagonaizer/effective-mobile/subscription-service/internal/config"
	"github.com/vagonaizer/effective-mobile/subscription-service/pkg/logger"
)

type Option func(*Server)

func WithLogger(log *logger.Logger) Option {
	return func(s *Server) {
		s.logger = log
	}
}

func WithConfig(cfg config.ServerConfig) Option {
	return func(s *Server) {
		s.config = cfg
	}
}

func WithRouter(router *gin.Engine) Option {
	return func(s *Server) {
		s.router = router
	}
}

func WithReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.writeTimeout = timeout
	}
}

func WithIdleTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.idleTimeout = timeout
	}
}

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}

func WithGracefulShutdown() Option {
	return func(s *Server) {
		s.enableGracefulShutdown = true
	}
}

func WithHealthCheck(healthCheckFunc func(ctx context.Context) error) Option {
	return func(s *Server) {
		s.healthCheck = healthCheckFunc
	}
}
