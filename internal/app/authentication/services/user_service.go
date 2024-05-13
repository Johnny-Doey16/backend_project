package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/utils"
	"github.com/steve-mir/diivix_backend/worker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MaxAccountRecoveryDuration is the duration for which an account can be recovered after deletion.
const (
	MaxAccountRecoveryDuration = 30 * 24 * time.Hour // for example, 30 days
	recoveryTokenLength        = 39
)

func DeleteAccountRequest(ctx context.Context, password string, store *sqlc.Store, uid uuid.UUID) error {

	// Check if the user exists and is not already marked as deleted
	user, err := store.GetUserByID(ctx, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return status.Error(codes.NotFound, "user with ID "+uid.String()+" not found")
		}
		return err
	}

	// Check password
	err = utils.CheckPassword(password, user.PasswordHash)
	if err != nil {
		return status.Error(codes.Unauthenticated, "password incorrect")
	}

	err = checkAccountStatus(user)
	if err != nil {
		return err
	}

	// Additional validations can be added here
	// ...

	// Proceed with soft deletion if all checks pass
	_, err = store.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:        uid,
		IsDeleted: sql.NullBool{Bool: true, Valid: true},
		DeletedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func AccRecoveryRequest(ctx context.Context, store *sqlc.Store, td worker.TaskDistributor, email string) error {
	// Check if the email is valid.
	if !utils.IsEmailFormat(email) {
		return fmt.Errorf("invalid email format: %s", email)
	}

	// Retrieve the user associated with the email address.
	user, err := store.GetUserByIdentifier(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("an error occurred while processing the recovery request")
		}
		return err
	}

	// Check if the user's account is already marked as deleted.
	if !user.IsDeleted.Bool { //!user.DeletedAt.Valid
		return errors.New("cannot process account recovery")
	}

	// Check if the account is within the recovery period.
	if time.Since(user.DeletedAt.Time) > MaxAccountRecoveryDuration {
		return fmt.Errorf("the account recovery period has expired for email: %s", email)
	}

	// Generate a recovery token and send an email to the user with the recovery instructions.
	// The token should be a secure, one-time use token with an expiry.
	// recoveryToken, err := generateSecureRecoveryToken()
	recoveryToken, err := utils.GenerateUniqueToken(recoveryTokenLength)
	if err != nil {
		return err
	}

	err = store.CreateUserDeleteRequest(ctx, sqlc.CreateUserDeleteRequestParams{
		UserID:        user.ID,
		Email:         user.Email,
		RecoveryToken: recoveryToken,
		ExpiresAt:     time.Now().Add(time.Minute * 15),
	})
	if err != nil {
		return err
	}

	err = SendEmail(td, ctx, email, recoveryToken)
	if err != nil {
		return err
	}

	return nil
}

func AccountRecovery(ctx context.Context, db *sql.DB, store *sqlc.Store, recoveryToken string) error {

	// Retrieve the user and recovery token information from the database
	usr, err := store.GetUserFromDeleteReqByToken(ctx, recoveryToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("account recovery request is invalid or has expired")
		}
		return err
	}

	if usr.Used.Bool {
		return errors.New("request token has been used or has expired")
	}

	if usr.ExpiresAt.Before(time.Now()) {
		return errors.New("account recovery request is invalid or has expired")
	}

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
	wg.Add(2)

	// Goroutine for updating token status
	go func() {
		defer wg.Done()

		// Assuming the recovery token is valid, proceed to unmark the account as deleted
		_, err = qtx.UpdateUser(ctx, sqlc.UpdateUserParams{
			ID:        usr.UserID,
			IsDeleted: sql.NullBool{Bool: false, Valid: true},
			DeletedAt: sql.NullTime{Time: time.Time{}, Valid: true}, // TODO: Set the null time
		})
		if err != nil {
			// Rollback the transaction in case of an error
			tx.Rollback()
		}
	}()

	// Goroutine for updating user account
	go func() {
		defer wg.Done()

		// Optionally, you may want to invalidate the recovery token after successful account recovery
		err = qtx.MarkDeleteAsUsedByToken(ctx, recoveryToken)
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
		return err
	}
	return nil
}
