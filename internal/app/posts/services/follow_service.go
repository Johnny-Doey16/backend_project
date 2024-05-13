package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/db/sqlc"
)

func Follow(ctx context.Context, db *sql.DB, store *sqlc.Store, follower, following uuid.UUID) error {
	if follower == following {
		return errors.New("cannot follow user")
	}

	// Create a context with a timeout for the transaction
	ctx, cancel := context.WithTimeout(ctx, time.Second*10) // Adjust the timeout as needed
	defer cancel()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	qtx := store.WithTx(tx)
	// Use a wait group to wait for both updates to complete
	var wg sync.WaitGroup
	wg.Add(3)

	// Insert into follow table
	// Goroutine for updating token status
	go func() {
		defer wg.Done()

		// Update password_request table
		err := addToFollowTable(ctx, qtx, follower, following)
		log.Println("Error 1:", err)
		if err != nil {
			tx.Rollback()
		}
	}()

	// Increase Follower in db
	// Goroutine for updating user password
	go func() {
		defer wg.Done()

		// Update users table
		err := increaseFollowers(ctx, qtx, follower)
		log.Println("Error 2:", err)
		if err != nil {
			// Rollback the transaction in case of an error
			tx.Rollback()
		}
	}()

	// Increase Following in db
	go func() {
		defer wg.Done()

		// Update users table
		err := increaseFollowings(ctx, qtx, following)
		log.Println("Error 3:", err)
		if err != nil {
			// Rollback the transaction in case of an error
			tx.Rollback()
		}
	}()

	// Wait for both goroutines to complete
	wg.Wait()

	// Commit the transaction if all updates were successful
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("could not follow user %s", err.Error())
	}

	return nil
}

func addToFollowTable(ctx context.Context, qtx *sqlc.Queries, follower, following uuid.UUID) error {
	err := qtx.CreateFollow(ctx, sqlc.CreateFollowParams{
		FollowerUserID:  follower,
		FollowingUserID: following,
	})
	if err != nil {
		return fmt.Errorf("cannot create follow data %s", err.Error())
	}
	return nil
}

func increaseFollowers(ctx context.Context, qtx *sqlc.Queries, follower uuid.UUID) error {
	err := qtx.UpdateIncreaseFollowers(ctx, follower)
	if err != nil {
		return fmt.Errorf("cannot increase followers %s", err.Error())
	}
	return nil
}

func increaseFollowings(ctx context.Context, qtx *sqlc.Queries, following uuid.UUID) error {
	err := qtx.UpdateIncreaseFollowing(ctx, following)
	if err != nil {
		return fmt.Errorf("cannot increase followers %s", err.Error())
	}
	return nil
}
