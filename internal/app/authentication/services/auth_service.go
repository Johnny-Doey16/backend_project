package services

import (
	"context"
	"database/sql"

	"github.com/steve-mir/diivix_backend/db/sqlc"
)

func BeginTrx(ctx context.Context, db *sql.DB, store *sqlc.Store) (*sql.Tx, *sqlc.Queries, error) {
	// Begin transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	qtx := store.WithTx(tx)

	return tx, qtx, nil
}
