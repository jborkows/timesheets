package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jborkows/timesheets/internal/model"
)

type TransactionSupport struct {
	queries *Queries
	db      *sql.DB
	config  *model.Config
}

func NewTransactionSupport(db *sql.DB) *TransactionSupport {
	dbTx := New(db)

	return &TransactionSupport{
		queries: dbTx,
		db:      db,
		config:  nil,
	}
}

func CreateRepository(db *sql.DB, config *model.Config) model.Repository {
	dbTx := New(db)

	return &TransactionSupport{
		queries: dbTx,
		db:      db,
		config:  config,
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
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %w (%v)", rbErr, err)
		}
		return fmt.Errorf("failed to run transaction: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil

}

func (support *TransactionSupport) Transactional(ctx context.Context, operation func(context.Context, model.Saver, model.Queryer) error) error {
	err := support.WithTransaction(ctx, func(ctx context.Context, q *Queries) error {
		repository := Repository(q, support.config.IsOvertime)
		return operation(ctx, repository, repository)
	})
	if err != nil {
		return fmt.Errorf("error in transaction: %w", err)
	}
	return nil
}
