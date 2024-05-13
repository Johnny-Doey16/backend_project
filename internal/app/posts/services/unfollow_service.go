package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnFollow(ctx context.Context, db *sql.DB, followerUid uuid.UUID, followingUid uuid.UUID) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return status.Errorf(codes.Internal, "transaction error: %v", err.Error())
	}
	defer tx.Rollback()

	// Executing the SQL code within the transaction
	query := fmt.Sprintf(`
	DO $$
    DECLARE
        rows_deleted INT;
    BEGIN
        DELETE FROM follow WHERE follower_user_id = '%s' AND following_user_id = '%s' RETURNING 1 INTO rows_deleted;
        IF rows_deleted > 0 THEN
            UPDATE entity_profiles SET following_count = following_count - 1 WHERE user_id = '%s';
			UPDATE entity_profiles SET followers_count = followers_count - 1 WHERE user_id = '%s';
        ELSE
            RAISE EXCEPTION 'No follow relationship found to delete';
        END IF;
    END $$;
    `, followerUid, followingUid, followerUid, followingUid)
	_, err = tx.Exec(query)

	if err != nil {
		// If there's an error, rollback the transaction
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func UnFollowOld(ctx context.Context, db *sql.DB, store *sqlc.Store, follower, following uuid.UUID) error {
	log.Println("Follower", follower)
	log.Println("Following", following)
	if follower == following {
		return errors.New("cannot unfollow user")
	}

	var err error
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
	// var wg sync.WaitGroup
	// wg.Add(3)

	var eg errgroup.Group

	// ! 1 Delete from Follow table
	eg.Go(func() error {
		err = qtx.UnFollow(ctx, sqlc.UnFollowParams{
			FollowerUserID:  follower,
			FollowingUserID: following,
		})
		log.Println("Deleting from follow table", err)
		if err != nil {
			return errors.New("an unknown error occurred deleting follow " + err.Error())
		}
		return nil
	})

	// ! 2
	eg.Go(func() error {
		err = qtx.DecreaseFollowers(ctx, follower)
		log.Println("decreasing followers", err)

		if err != nil {
			return errors.New("an unknown error occurred decreasing followers " + err.Error())
		}
		return nil
	})

	// ! 3. Works
	eg.Go(func() error {
		err = qtx.DecreaseFollowing(ctx, following)
		log.Println("decreasing followings", err)

		if err != nil {
			return errors.New("an unknown error occurred decreasing following " + err.Error())
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction if both goroutines completed successfully
	return tx.Commit()

	// Insert into follow table
	// Goroutine for updating token status
	// go func() {
	// 	defer wg.Done()

	// 	// Update password_request table
	// 	err := deleteFromFollowTable(ctx, qtx, follower, following)
	// 	if err != nil {
	// 		tx.Rollback()
	// 	}
	// }()

	// // Increase Follower in db
	// // Goroutine for updating user password
	// go func() {
	// 	defer wg.Done()

	// 	// Update users table
	// 	err := decreaseFollowers(ctx, qtx, follower)
	// 	if err != nil {
	// 		// Rollback the transaction in case of an error
	// 		tx.Rollback()
	// 	}
	// }()

	// // Increase Following in db
	// go func() {
	// 	defer wg.Done()

	// 	// Update users table
	// 	err := decreaseFollowings(ctx, qtx, following)
	// 	if err != nil {
	// 		// Rollback the transaction in case of an error
	// 		tx.Rollback()
	// 	}
	// }()

	// // Wait for both goroutines to complete
	// wg.Wait()

	// // Commit the transaction if all updates were successful
	// err = tx.Commit()
	// if err != nil {
	// 	return fmt.Errorf("could not follow user %s", err.Error())
	// }

	// return nil
}
