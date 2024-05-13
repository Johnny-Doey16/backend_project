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
	e "github.com/steve-mir/diivix_backend/errors"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	au "github.com/steve-mir/diivix_backend/internal/app/authentication/utils"
	"github.com/steve-mir/diivix_backend/utils"
	"github.com/steve-mir/diivix_backend/worker"
)

// Request change
func RequestChangeOfEmail(taskDistributor worker.TaskDistributor, ctx context.Context, store *sqlc.Store, uid uuid.UUID, email, newEmail string) error {
	// Validate the new email format
	if err := au.ValidateEmail(newEmail); err != nil {
		return fmt.Errorf("%s: %s", e.ErrInvalidEmailFormat, err.Error())
	}

	// Initialize the change request
	_, resetCode, err := initChangeRequest(ctx, store, email)
	if err != nil {
		return fmt.Errorf("failed to initialize change request %v", err.Error())
	}

	// Construct the message for email
	msg := fmt.Sprintf("Below is code to change your email: %s.\nPlease do not share this with anyone", resetCode)
	SendEmail(taskDistributor, ctx, newEmail, msg)

	// Add link to the database
	err = store.CreateChangeIdRequest(ctx, sqlc.CreateChangeIdRequestParams{
		UserID:     uid,
		Identifier: newEmail,
		Type:       "email",
		Token:      resetCode,
		ExpiresAt:  time.Now().Add(time.Minute * 15),
	})
	if err != nil {
		return fmt.Errorf("failed to create change ID request %v", err.Error())
	}

	return nil
}

func RequestChangeOfPhone(ctx context.Context, store *sqlc.Store, uid uuid.UUID, uPhone, newPhone string) error {
	if au.ValidatePhone(newPhone) {
		return fmt.Errorf(e.InvalidPhoneFormat)
	}

	usr, resetCode, err := initChangeRequest(ctx, store, uPhone)
	if err != nil {
		return errors.New(UnexpectedErr + " 3 " + err.Error())
	}
	msg := fmt.Sprintf("Below is code to change your phone: %s.\nPlease do not share this with anyone", resetCode)

	// TODO send this to their phone
	log.Printf("Sending MSG: %s. TO: %s", msg, usr.Phone.String)

	// Add link to db
	err = store.CreateChangeIdRequest(ctx, sqlc.CreateChangeIdRequestParams{
		UserID:     uid,
		Identifier: newPhone,
		Type:       "phone",
		Token:      resetCode,
		ExpiresAt:  time.Now().Add(time.Minute * 15),
	})
	if err != nil {
		return errors.New(UnexpectedErr + " " + err.Error())
	}

	return nil
}

// Confirm change
func ChangeEmail(ctx context.Context, db *sql.DB, config utils.Config, store *sqlc.Store, claims *token.Payload, code string) (string, time.Time, error) {
	// Create a context with a timeout for the transaction
	ctx, cancel := context.WithTimeout(ctx, time.Second*10) // Adjust the timeout as needed
	defer cancel()

	if len(code) != length {
		return "", time.Time{}, errors.New("invalid token")
	}

	tokenData, err := store.GetChangeIdRequestByToken(ctx, code)
	if err != nil {
		return "", time.Time{}, err
	}

	if tokenData.UserID != claims.UserId {
		return "", time.Time{}, errors.New("cannot change identifier")
	}

	if tokenData.ExpiresAt.Before(time.Now()) {
		return "", time.Time{}, fmt.Errorf("token expired")
	}

	if tokenData.Used.Bool {
		return "", time.Time{}, fmt.Errorf("token already used")
	}

	// Generate new accessToken
	accessToken, accessExpires, err := generateNewTokens(tokenData.Identifier, claims.Username, claims.Phone, claims, config)
	if err != nil {
		return "", time.Time{}, err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return "", time.Time{}, err
	}
	defer tx.Rollback()

	/*
		qtx := store.WithTx(tx)
		// Use a wait group to wait for both updates to complete
		var wg sync.WaitGroup
		wg.Add(2)

		// ! Goroutine for updating token status
		go func() {
			defer wg.Done()

			// Update password_request table
			err := updateTokenStatus(ctx, qtx, code)
			if err != nil {
				fmt.Println("Token stat", err)
				// Rollback the transaction in case of an error
				tx.Rollback()
			}
		}()

		// Goroutine for updating user email or phone number
		go func() {
			defer wg.Done()

			// Update email users table
			err := updateUserEmail(ctx, qtx, tokenData.UserID, tokenData.Identifier)
			if err != nil {
				fmt.Println("Update user error", err)
				// Rollback the transaction in case of an error
				tx.Rollback()
			}

		}()

		// Wait for both goroutines to complete
		wg.Wait()

		// Commit the transaction if all updates were successful
		err = tx.Commit()
		if err != nil {
			return "", time.Time{}, err
		}
	*/

	//*Using chans (TEST)
	qtx := store.WithTx(tx)

	// Use channels for goroutine errors
	errCh := make(chan error, 2)

	// Goroutine for updating token status
	go func() {
		errCh <- updateTokenStatus(ctx, qtx, code)
	}()

	// Goroutine for updating user email
	go func() {
		errCh <- updateUserEmail(ctx, qtx, tokenData.UserID, tokenData.Identifier)
	}()

	// Wait for goroutines to complete and collect errors
	for i := 0; i < 2; i++ {
		if err := <-errCh; err != nil {
			tx.Rollback()
			return "", time.Time{}, err
		}
	}

	// Commit the transaction if all updates were successful
	if err := tx.Commit(); err != nil {
		return "", time.Time{}, err
	}

	return accessToken, accessExpires, nil

}

func ChangePhone(ctx context.Context, db *sql.DB, config utils.Config, store *sqlc.Store, claims *token.Payload, code string) (string, time.Time, error) {
	// Create a context with a timeout for the transaction
	ctx, cancel := context.WithTimeout(ctx, time.Second*10) // Adjust the timeout as needed
	defer cancel()

	if len(code) != length {
		return "", time.Time{}, errors.New("invalid token")
	}

	tokenData, err := store.GetChangeIdRequestByToken(ctx, code)
	if err != nil {
		return "", time.Time{}, err
	}

	if tokenData.UserID != claims.UserId {
		return "", time.Time{}, errors.New("cannot change identifier")
	}

	if tokenData.ExpiresAt.Before(time.Now()) {
		return "", time.Time{}, fmt.Errorf("token expired")
	}

	if tokenData.Used.Bool {
		return "", time.Time{}, fmt.Errorf("token already used")
	}

	// Generate new accessToken
	accessToken, accessExpires, err := generateNewTokens(claims.Email, claims.Username, tokenData.Identifier, claims, config)
	if err != nil {
		return "", time.Time{}, err
	}

	tx, err := db.Begin()
	if err != nil {
		return "", time.Time{}, err
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
		err := updateTokenStatus(ctx, qtx, code)
		if err != nil {
			// Rollback the transaction in case of an error
			tx.Rollback()
		}
	}()

	// Goroutine for updating user phone number
	go func() {
		defer wg.Done()

		// Update the phone users table
		err := updateUserPhone(ctx, qtx, tokenData.UserID, tokenData.Identifier)
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
		return "", time.Time{}, err
	}

	return accessToken, accessExpires, nil

}

// *** Change Username *** *** *** ***
func ChangeUsername(ctx context.Context, db *sql.DB, config utils.Config, store *sqlc.Store, claims *token.Payload, username string) (string, time.Time, error) {
	// Verify claims (optional)

	// Verify newUsername format
	if !au.ValidateUsername(username) {
		return "", time.Time{}, fmt.Errorf(e.InvalidUsername)
	}

	//* Use transaction to update username
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return "", time.Time{}, err
	}
	defer tx.Rollback()

	qtx := store.WithTx(tx)

	err = CheckIfUsernameExists(ctx, qtx, username)
	if err != nil {
		log.Println("Error getting username:", err)
		tx.Rollback()
	}

	err = updateUsername(ctx, qtx, claims.UserId, username)
	if err != nil {
		log.Println("Error updating username:", err)
		tx.Rollback()
	}

	// Generate new accessToken
	accessToken, accessExpires, err := generateNewTokens(claims.Email, username, claims.Phone, claims, config)
	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken, accessExpires, tx.Commit()

}

// ***
func updateUserEmail(ctx context.Context, qtx *sqlc.Queries, id uuid.UUID, email string) error {
	_, err := qtx.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:    id,
		Email: sql.NullString{String: email, Valid: true},
	})
	return err
}

func updateUserPhone(ctx context.Context, qtx *sqlc.Queries, id uuid.UUID, phone string) error {
	_, err := qtx.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:    id,
		Phone: sql.NullString{String: phone, Valid: true},
	})
	return err
}

func updateUsername(ctx context.Context, qtx *sqlc.Queries, id uuid.UUID, username string) error {
	_, err := qtx.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:       id,
		Username: sql.NullString{String: username, Valid: true},
	})
	return err
}

func CheckIfUsernameExists(ctx context.Context, qtx *sqlc.Queries, username string) error {
	// Compares the username lowercase with lowercase of that in db
	c, err := qtx.CheckUsername(ctx, username)
	if err != nil {
		return err
	}

	// If count is greater than 0 it means a username variant was found irrespective of the case
	if c > 0 {
		return errors.New("username taken")
	}

	// If none of the conditions are met, it means no user exists with that username and there were no errors.
	return nil
}

func generateNewTokens(email, username, phone string, claims *token.Payload, config utils.Config) (string, time.Time, error) {
	// Generate tokens with new data
	tokenService := NewTokenService(config)

	// Access token
	accessToken, accessPayload, err := tokenService.CreateAccessToken(email, username, phone, claims.MfaPassed, claims.IsEmailVerified, claims.UserId, claims.Role, claims.IP, claims.UserAgent)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error creating access token %s", err)
	}

	return accessToken, accessPayload.Expires, nil
}
