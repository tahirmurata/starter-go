package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"starter/internal/config"
	"starter/internal/logger"
	"starter/internal/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ServiceManager represents a service that interacts with a database.
type ServiceManager interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health(ctx context.Context) (stats map[string]string, err error)

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close()
}

// Manager represents a manager for database operations.
type Manager interface {
	// InWriteTx runs the passed in function within a transaction with the given isolation level.
	InWriteTx(ctx context.Context, level pgx.TxIsoLevel, fn func(tx pgx.Tx, q *sqlc.Queries) error) error

	// InReadTx runs the passed in function within a read-only transaction with a read commited isolation level.
	InReadTx(ctx context.Context, fn func(tx pgx.Tx, q *sqlc.Queries) error) error

	// Close closes the database connection pool.
	Close()
}

type database struct {
	pool *pgxpool.Pool
}

type service struct {
	db Manager
}

var dbInstance *service

func New(cfg *config.Config) ServiceManager {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	db, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database, cfg.Database.Schema))
	if err != nil {
		slog.LogAttrs(context.Background(), logger.LevelFatal, "Failed to create database pool", slog.Any("err", err))
		panic(1)
	}

	dbInstance = &service{
		db: &database{
			pool: db,
		},
	}

	return dbInstance
}

func (s *service) Health(ctx context.Context) (stats map[string]string, err error) {
	if err := s.db.InReadTx(ctx, func(tx pgx.Tx, q *sqlc.Queries) error {
		// Ping the database
		if err := tx.Conn().Ping(ctx); err != nil {
			stats["status"] = "down"
			stats["error"] = fmt.Sprintf("db down: %v", err)
			return err
		}

		// Database is up, add more statistics
		stats["status"] = "up"
		stats["message"] = "bing chilling"
		return nil
	}); err != nil {
		return nil, errors.Join(errors.New("failed to ping database"), err)
	}

	return stats, nil
}

func (s *service) Close() {
	s.db.Close()
}
