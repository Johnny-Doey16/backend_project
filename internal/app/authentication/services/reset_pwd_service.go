package services

import (
	"context"
	"database/sql"

	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	e "github.com/pkg/errors"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/utils"
	"github.com/steve-mir/diivix_backend/worker"
)

const (
	codeLength    = int64(1000000)
	length        = 6
	ResetMsg      = "if an account exists a password reset email will be sent to you"
	UnexpectedErr = "an unexpected error occurred"
)

func RequestPwdReset(ctx context.Context, store *sqlc.Store, td worker.TaskDistributor, email string) error {
	usr, pwdResetCodeStr, err := initChangeRequest(ctx, store, email)
	if err != nil {
		return errors.New(ResetMsg)
	}
	msg := fmt.Sprintf("Below is code to reset your password: %s.\nPlease do not share this with anyone", pwdResetCodeStr)

	errChan := make(chan error, 2) // Buffer the channel to the number of goroutines
	var wg sync.WaitGroup

	// Add link to db
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := store.CreatePasswordResetRequest(ctx, sqlc.CreatePasswordResetRequestParams{
			UserID:    usr.ID,
			Email:     usr.Email,
			Token:     pwdResetCodeStr,
			ExpiresAt: time.Now().Add(time.Minute * 15),
		}); err != nil {
			errChan <- err
		}
	}()

	// Send email
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := SendEmail(td, ctx, email, msg); err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return fmt.Errorf("an unexpected error occurred: %v", err)
		}
	}

	return nil
}

func ResetPassword(ctx context.Context, db *sql.DB, store *sqlc.Store, code, pwd string) error {
	// Create a context with a timeout for the transaction
	ctx, cancel := context.WithTimeout(ctx, time.Second*10) // Adjust the timeout as needed
	defer cancel()

	if len(code) != length {
		return errors.New("invalid token")
	}

	if !utils.ValidatePassword(pwd) {
		return errors.New("invalid password format")
	}

	tokenData, err := store.GetPasswordResetRequestByToken(ctx, code)
	if err != nil {
		return err
	}

	if tokenData.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("token expired")
	}

	if tokenData.Used.Bool {
		return fmt.Errorf("token already used")
	}

	// ! 1 Get User
	user, err := getUser(ctx, store, tokenData.Email)
	if err != nil {
		return err
	}

	err = checkAccountStatus(user)
	if err != nil {
		return err
	}

	// ! 2 Check old password
	if err = utils.CheckPassword(pwd, user.PasswordHash); err == nil {
		return errors.New("cannot use old password")
	}

	// ! 3 Hash password
	hashedPwd, err := utils.HashPassword(pwd)
	if err != nil {
		return err
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

		// Update password_request table
		err := updateResetPwdTokenStatus(ctx, qtx, code)
		if err != nil {
			// Rollback the transaction in case of an error
			tx.Rollback()
		}
	}()

	// Goroutine for updating user password
	go func() {
		defer wg.Done()

		// Update users table
		err := updateUserPassword(ctx, qtx, tokenData.Email, hashedPwd)
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

func getUser(ctx context.Context, store *sqlc.Store, email string) (sqlc.Authentication, error) {
	user, err := store.GetUserByIdentifier(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return sqlc.Authentication{}, errors.New("1. email not found")
		}
		return sqlc.Authentication{}, errors.New("email not found")
	}
	return user, nil
}

func checkAccountStatus(usr sqlc.Authentication) error {
	if usr.IsSuspended.Bool {
		return errors.New("account suspended")
	}

	if usr.IsDeleted.Bool {
		return errors.New("account deleted")
	}
	return nil
}

func updateResetPwdTokenStatus(ctx context.Context, qtx *sqlc.Queries, token string) error {
	log.Println("Inside update token")
	// ! 4 Use transaction to update these
	err := qtx.UpdatePasswordResetRequestByToken(ctx, sqlc.UpdatePasswordResetRequestByTokenParams{
		Token: token,
		Used:  sql.NullBool{Bool: true, Valid: true},
	})
	log.Println("Done with update token")
	return err
}

func updateTokenStatus(ctx context.Context, qtx *sqlc.Queries, token string) error {
	log.Println("Inside update token")
	// ! 4 Use transaction to update these
	err := qtx.UpdateChangeIdByToken(ctx, sqlc.UpdateChangeIdByTokenParams{
		Token: token,
		Used:  sql.NullBool{Bool: true, Valid: true},
	})
	log.Println("Done with update token")
	return err
}

func updateUserPassword(ctx context.Context, qtx *sqlc.Queries, email, hashedPwd string) error {
	// Change password
	err := qtx.UpdateUserPasswordByEmail(ctx, sqlc.UpdateUserPasswordByEmailParams{
		Email:        email,
		PasswordHash: hashedPwd,
	})
	return err
}

func initChangeRequest(ctx context.Context, store *sqlc.Store, id string) (sqlc.Authentication, string, error) {
	emailResetCode, err := utils.GenerateSecureRandomNumber(codeLength)
	if err != nil {
		return sqlc.Authentication{}, "", e.Wrap(err, "failed to generate secure random number")
	}

	// Check if identifier exists
	usr, err := store.GetUserByIdentifier(ctx, id)
	if err != nil {
		return sqlc.Authentication{}, "", e.Wrap(err, "failed to get user by identifier. "+ResetMsg)
	}

	// check account status
	err = checkAccountStatus(usr)
	if err != nil {
		return sqlc.Authentication{}, "", err
	}

	return usr, fmt.Sprintf("%06d", emailResetCode), nil
}
