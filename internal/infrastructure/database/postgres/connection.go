package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/vagonaizer/effective-mobile/subscription-service/internal/config"
	"github.com/vagonaizer/effective-mobile/subscription-service/pkg/logger"
)

type DB struct {
	pool *pgxpool.Pool
	log  *logger.Logger
}

func New(cfg config.DatabaseConfig, log *logger.Logger) (*DB, error) {
	log.Info("connecting to postgres",
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port),
		zap.String("database", cfg.DBName))

	poolConfig, err := buildPoolConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("build pool config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	db := &DB{
		pool: pool,
		log:  log,
	}

	if err := db.ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	log.Info("postgres connected successfully",
		zap.Int32("max_conns", poolConfig.MaxConns),
		zap.Int32("min_conns", poolConfig.MinConns))

	return db, nil
}

func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

func (db *DB) Close() {
	if db.pool != nil {
		db.pool.Close()
		db.log.Info("postgres connection closed")
	}
}

func (db *DB) ping(ctx context.Context) error {
	if err := db.pool.Ping(ctx); err != nil {
		db.log.Error("postgres ping failed", zap.Error(err))
		return fmt.Errorf("ping database: %w", err)
	}
	return nil
}

func (db *DB) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return db.ping(ctx)
}

func (db *DB) Stats() *pgxpool.Stat {
	return db.pool.Stat()
}

func buildPoolConfig(cfg config.DatabaseConfig) (*pgxpool.Config, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.MaxIdleConns)
	poolConfig.MaxConnLifetime = time.Duration(cfg.MaxLifetime) * time.Second
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	return poolConfig, nil
}
