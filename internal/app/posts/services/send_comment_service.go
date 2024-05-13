package services

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type commentStruct struct {
	commentID int32
	err       error
}

type commentCountStruct struct {
	totalComments int32
	err           error
}

func CheckComment(content string) error {
	if content == "" {
		return fmt.Errorf("body cannot be empty")
	}

	if len(content) < minContent {
		return fmt.Errorf("characters must be at least %v characters long", minContent)
	}

	if len(content) > maxContent {
		return fmt.Errorf("characters must not be more than %v characters long", maxContent)
	}

	return nil
}

func RunPostCommentConcurrent(tx *sql.Tx, qtx *sqlc.Queries, ctx context.Context, content string, postID, uid uuid.UUID) (int32, int32, error) {
	var wg sync.WaitGroup
	commentChan := make(chan commentStruct, 1)
	incrementCommentChan := make(chan commentCountStruct, 1)

	wg.Add(1) // We have 2 Goroutines
	go func() {
		defer wg.Done()
		comment, err := qtx.CreateComment(ctx, sqlc.CreateCommentParams{
			PostID:      postID,
			CommentText: content,
			UserID:      uid,
		})
		// You must check whether the operation was successful before sending the ID
		var id int32
		if err == nil {
			id = comment.ID
		}
		commentChan <- commentStruct{commentID: id, err: err}
	}()
	wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		c_counts, err := qtx.IncrementComments(ctx, postID)
		var count int32
		if err == nil {
			count = c_counts.Comments.Int32
		}
		incrementCommentChan <- commentCountStruct{totalComments: count, err: err}
	}()

	wg.Wait() // Wait for both Goroutines to finish
	close(commentChan)
	close(incrementCommentChan)

	comment, counts := <-commentChan, <-incrementCommentChan

	// Check for errors and rollback if necessary
	if comment.err != nil || counts.err != nil {
		tx.Rollback()
		if comment.err != nil {
			return 0, 0, status.Errorf(codes.Aborted, "an unknown error occurred creating comment: %v", comment.err)
		}
		if counts.err != nil {
			return 0, 0, status.Errorf(codes.Aborted, "an unknown error occurred incrementing comment count: %v", counts.err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, status.Errorf(codes.Aborted, "an unknown error occurred committing the transaction: %v", err)
	}

	return comment.commentID, counts.totalComments, nil
}

type CommentUpdate struct {
	PostID string `json:"post_id"`
	UserID string `json:"user_id"`
}
