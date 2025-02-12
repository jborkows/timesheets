package db

import (
	"context"
	"database/sql"
	"fmt"
)

type TransactionSupport struct {
	queries *Queries
	db      *sql.DB
}

func NewTransactionSupport(db *sql.DB) *TransactionSupport {
	dbTx := New(db)

	return &TransactionSupport{
		queries: dbTx,
		db:      db,
	}
}

func (support *TransactionSupport) WithTransaction(ctx context.Context, operation func(context.Context, *Queries) error) error {
	tx, err := support.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	workUnit := support.queries.WithTx(tx)
	err = operation(ctx, workUnit)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("failed to rollback transaction: %w", err)
		}
		return fmt.Errorf("failed to run transaction: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil

}
