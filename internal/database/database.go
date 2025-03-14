package database

import (
	"context"
	"fmt"
	"starter/internal/config"
	"starter/internal/sqlc"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Database represents a service that interacts with a database.
type Database interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health(ctx context.Context) (stats map[string]string, err error)

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close()
}

// helper represents a manager for database operations.
type helper interface {
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
	db helper
}

var dbInstance *service
var once sync.Once

func New(ctx context.Context, cfg *config.Config) (Database, error) {
	var initError error
	once.Do(func() {
		db, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database, cfg.Database.Schema))
		if err != nil {
			initError = fmt.Errorf("creating pgx pool: %w", err)
			return
		}

		dbInstance = &service{
			db: &database{
				pool: db,
			},
		}
	})

	return dbInstance, initError
}

func (s *service) Health(ctx context.Context) (map[string]string, error) {
	stats := make(map[string]string)
	err := s.db.InReadTx(ctx, func(tx pgx.Tx, q *sqlc.Queries) error {
		// Ping the database
		if err := tx.Conn().Ping(ctx); err != nil {
			stats["status"] = "down"
			stats["error"] = fmt.Sprintf("db down: %v", err)
			return fmt.Errorf("ping database: %w", err)
		}

		// Database is up, add more statistics
		stats["status"] = "up"
		stats["message"] = "bing chilling"
		return nil
	})

	return stats, err
}

func (s *service) Close() {
	s.db.Close()
}
