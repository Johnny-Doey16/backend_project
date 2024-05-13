package services

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/db/sqlc"
)

type ViewUpdate struct {
	PostID string `json:"postId"`
	UserID string `json:"userId"`
}

// TODO: And views table to avoid duplicate viewing
func ViewPostInDbConcurrently(tx *sql.Tx, qtx *sqlc.Queries, ctx context.Context, postID, uid uuid.UUID) error {

	var wg sync.WaitGroup
	incrementViewChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := qtx.IncrementViewCount(ctx, postID)
		incrementViewChan <- err
	}()

	wg.Wait()
	close(incrementViewChan)

	if err := <-incrementViewChan; err != nil {
		tx.Rollback()
		return errors.New("an unknown error occurred incrementing like " + err.Error())
	}

	return tx.Commit()
}
