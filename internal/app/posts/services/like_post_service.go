package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/db/sqlc"
)

type LikeUpdate struct {
	PostID string `json:"postId"`
	UserID string `json:"userId"`
}

// TODO: Try optimizing with go routine
func LikePostInDbConcurrently(tx *sql.Tx, qtx *sqlc.Queries, ctx context.Context, postID, uid uuid.UUID) (int32, bool, error) {
	ls, err := qtx.LikePost(ctx, sqlc.LikePostParams{
		UserID: uid,
		PostID: postID,
	})
	if err != nil {
		tx.Rollback()
		return 0, false, errors.New("an unknown error occurred crAdding like " + err.Error())
	}

	return ls[0].Likes.Int32, ls[0].Liked, tx.Commit()
}

// * Depreciated functions
func UnLikePostInDbConcurrently(tx *sql.Tx, qtx *sqlc.Queries, ctx context.Context, postID, uid uuid.UUID) (int32, error) {
	err := qtx.RemoveLike(ctx, sqlc.RemoveLikeParams{
		UserID: uid,
		PostID: postID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			// No like to remove, not necessarily an error, handle according to your use case
			fmt.Println("No such row")
		} else {
			if err != nil {
				tx.Rollback()
				return 0, fmt.Errorf("an unknown error occurred removing like: %w", err)
			}
		}
	} else {
		// Like removed successfully, handle the id if needed
		fmt.Println("Like removed successfully")
	}

	lD, err := qtx.DecrementLikeCount(ctx, postID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("an unknown error occurred decrementing like: %w", err)
	}

	return lD.Likes.Int32, tx.Commit()
}

/*
func LikePostInDbConcurrently(tx *sql.Tx, qtx *sqlc.Queries, ctx context.Context, postID, uid uuid.UUID) (int32, error) {
	var pM sqlc.PostsMetric
	var ls int32
	var eg errgroup.Group
	var err error

	// ! 1
	eg.Go(func() error {
		err = qtx.AddLike(ctx, sqlc.AddLikeParams{
			UserID: uid,
			PostID: postID,
		})
		if err != nil {
			return errors.New("an unknown error occurred crAdding like " + err.Error())
		}
		return nil
	})

	// ! 2
	eg.Go(func() error {
		pM, err = qtx.IncrementLikeCount(ctx, postID)

		if err != nil {
			return errors.New("an unknown error occurred incrementing like " + err.Error())
		}
		ls = pM.Likes.Int32
		return nil
	})

	if err := eg.Wait(); err != nil {
		tx.Rollback()
		return 0, err
	}

	// Commit the transaction if both goroutines completed successfully
	return ls, tx.Commit()
}

func UnLikePostInDbConcurrently(tx *sql.Tx, qtx *sqlc.Queries, ctx context.Context, postID, uid uuid.UUID) (int32, error) {
	var pM sqlc.PostsMetric
	var ls int32
	var eg errgroup.Group
	var err error

	// ! 1
	eg.Go(func() error {
		err = qtx.RemoveLike(ctx, sqlc.RemoveLikeParams{
			UserID: uid,
			PostID: postID,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				// No like to remove, not necessarily an error, handle according to your use case
				fmt.Println("No such row")
				// return nil
			} else {
				if err != nil {
					// tx.Rollback()
					return fmt.Errorf("an unknown error occurred removing like: %w", err)
				}
			}
		} else {
			// Like removed successfully, handle the id if needed
			fmt.Println("Like removed successfully")
			// return nil
		}

		return nil
	})

	// ! 2
	eg.Go(func() error {
		pM, err = qtx.DecrementLikeCount(ctx, postID)

		if err != nil {
			return fmt.Errorf("an unknown error occurred decrementing like: %w", err)
		}
		ls = pM.Likes.Int32
		return nil
	})

	if err := eg.Wait(); err != nil {
		tx.Rollback()
		return 0, err
	}

	return ls, tx.Commit()
}
*/
