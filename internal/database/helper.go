package database

import (
	"context"
	"errors"
	"fmt"
	"starter/internal/sqlc"

	"github.com/jackc/pgx/v5"
)

func (db *database) InWriteTx(ctx context.Context, level pgx.TxIsoLevel, fn func(tx pgx.Tx, q *sqlc.Queries) error) error {
	return db.inTx(ctx, level, pgx.ReadWrite, pgx.NotDeferrable, fn)
}

func (db *database) InReadTx(ctx context.Context, fn func(tx pgx.Tx, q *sqlc.Queries) error) error {
	return db.inTx(ctx, pgx.RepeatableRead, pgx.ReadOnly, pgx.NotDeferrable, fn)
}

// inTx runs the passed in function within a transaction with the given transaction options.
func (db *database) inTx(ctx context.Context, level pgx.TxIsoLevel, access pgx.TxAccessMode, deferrable pgx.TxDeferrableMode, fn func(tx pgx.Tx, q *sqlc.Queries) error) (err error) {
	conn, errAcquire := db.pool.Acquire(ctx)
	if errAcquire != nil {
		return fmt.Errorf("acquire connection: %w", errAcquire)
	}
	defer conn.Release()

	opts := pgx.TxOptions{
		IsoLevel:       level,
		AccessMode:     access,
		DeferrableMode: deferrable,
	}

	tx, errBegin := conn.BeginTx(ctx, opts)
	if errBegin != nil {
		return fmt.Errorf("begin tx: %w", errBegin)
	}

	defer func() {
		errRollback := tx.Rollback(ctx)
		if !(errRollback == nil || errors.Is(errRollback, pgx.ErrTxClosed)) {
			err = errors.Join(err, errRollback)
		}
	}()

	if err := fn(tx, sqlc.New(tx)); err != nil {
		return fmt.Errorf("run function: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (db *database) Close() {
	db.pool.Close()
}
