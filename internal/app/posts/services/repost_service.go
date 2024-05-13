package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/cache"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type repostStruct struct {
	repost sqlc.Repost
	err    error
}

type repostCountStruct struct {
	totalReposts int32
	err          error
}

func RunPostRepostConcurrent(tx *sql.Tx, qtx *sqlc.Queries, ctx context.Context, originalPostID, uid uuid.UUID) (sqlc.Repost, int32, error) {
	var wg sync.WaitGroup
	repostChan := make(chan repostStruct, 1)
	incrementRepostsChan := make(chan repostCountStruct, 1)

	wg.Add(1) // We have 2 Goroutines
	go func() {
		defer wg.Done()
		repost, err := qtx.CreateRepost(ctx, sqlc.CreateRepostParams{
			UserID:         uid,
			OriginalPostID: originalPostID,
		})
		repostChan <- repostStruct{repost: repost, err: err}
	}()
	wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		r_counts, err := qtx.IncrementReposts(ctx, originalPostID)
		var count int32
		if err == nil {
			count = r_counts.Int32
		}
		incrementRepostsChan <- repostCountStruct{totalReposts: count, err: err}
	}()

	wg.Wait() // Wait for both Goroutines to finish
	close(repostChan)
	close(incrementRepostsChan)

	repost, counts := <-repostChan, <-incrementRepostsChan

	// Check for errors and rollback if necessary
	if repost.err != nil || counts.err != nil {
		tx.Rollback()
		if repost.err != nil {
			return sqlc.Repost{}, 0, status.Errorf(codes.Aborted, "an unknown error occurred creating repost: %v", repost.err)
		}
		if counts.err != nil {
			return sqlc.Repost{}, 0, status.Errorf(codes.Aborted, "an unknown error occurred incrementing repost count: %v", counts.err)
		}
	}

	if err := tx.Commit(); err != nil {
		return sqlc.Repost{}, 0, status.Errorf(codes.Aborted, "an unknown error occurred committing the transaction: %v", err)
	}

	return repost.repost, counts.totalReposts, nil
}

type RepostUpdate struct {
	PostID string `json:"post_id"`
	UserID string `json:"user_id"`
}

func AddRepostCountToRedis(redisCache cache.Cache, ctx context.Context, key, postID string, counts int32) error {
	// Set the repost count in Redis to match the database count
	commentCountKey := fmt.Sprintf("%s:%s", key, postID)
	err := redisCache.SetKey(ctx, commentCountKey, counts, time.Hour*24)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to set repost count in Redis: %v", err)
	}
	return nil
}

func BroadcastRepostEvent(ctx context.Context, redisCache cache.Cache,
	postId, repostChanKey string, uid uuid.UUID) error {
	repostUpdate := RepostUpdate{
		PostID: postId,
		UserID: uid.String(),
		// Additional fields as needed
	}
	repostUpdateBytes, err := json.Marshal(repostUpdate)
	if err != nil {
		log.Printf("Error marshalling comment update: %v", err)
		return status.Errorf(codes.Internal, "Failed to marshalling repost: %v", err)
	}

	// Publish the like update to the Redis channel
	err = redisCache.Publish(ctx, repostChanKey, repostUpdateBytes).Err()
	if err != nil {
		log.Printf("Error publishing like update: %v", err)
		return status.Errorf(codes.Internal, "Failed to publishing repost to redis: %v", err)
	}
	return nil
}
