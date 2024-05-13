package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	e "github.com/pkg/errors"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/utils"
	"github.com/steve-mir/diivix_backend/worker"
)

func SendVerificationEmail(qtx *sqlc.Queries, ctx context.Context, td worker.TaskDistributor, userId uuid.UUID, email string) error {

	verificationCode, err := utils.GenerateSecureRandomNumber(codeLength)
	if err != nil {
		return e.Wrap(err, "failed to generate secure random number")
	}
	code := fmt.Sprintf("%06d", verificationCode)
	content := fmt.Sprintf("Use this to verify you email. code %s", code)

	// Use a WaitGroup to wait for both goroutines to complete
	var wg sync.WaitGroup
	wg.Add(2) // We have two goroutines

	// Error channel with buffer for two errors
	errChan := make(chan error, 2)

	go func() {
		defer wg.Done() // Notify the WaitGroup that this goroutine is done
		// Send email here.
		if err := SendEmail(td, ctx, email, content); err != nil {
			errChan <- e.Wrap(err, "failed to send verification email")
		}
	}()

	go func() {
		defer wg.Done() // Notify the WaitGroup that this goroutine is done
		// Add link to db
		if err := qtx.CreateEmailVerificationRequest(ctx, sqlc.CreateEmailVerificationRequestParams{
			UserID:    userId,
			Email:     email,
			Token:     code,
			ExpiresAt: time.Now().Add(time.Minute * 15),
		}); err != nil {
			errChan <- e.Wrap(err, "failed to create email verification request")
		}
	}()

	// Wait for both goroutines to complete
	wg.Wait()
	close(errChan) // Close the channel so that the range loop can finish

	// Collect errors from the error channel
	for err := range errChan {
		if err != nil {
			return err // Return the first error encountered
		}
	}

	return nil

}

func ReSendVerificationEmail(store *sqlc.Store, ctx context.Context, td worker.TaskDistributor, userId uuid.UUID, email string) error {

	// Check if identifier exists
	usr, err := store.GetUserByID(ctx, userId)
	if err != nil {
		return e.Wrap(err, "failed to get user by identifier. "+ResetMsg)
	}

	// check account status
	err = checkAccountStatusForEmail(usr)
	if err != nil {
		return err
	}

	verificationCode, err := utils.GenerateSecureRandomNumber(codeLength)
	if err != nil {
		return e.Wrap(err, "failed to generate secure random number")
	}
	code := fmt.Sprintf("%06d", verificationCode)
	content := fmt.Sprintf("Use this to verify you email. code %s", code)

	// Use a WaitGroup to wait for both goroutines to complete
	var wg sync.WaitGroup
	wg.Add(2) // We have two goroutines

	// Error channel with buffer for two errors
	errChan := make(chan error, 2)

	go func() {
		defer wg.Done() // Notify the WaitGroup that this goroutine is done
		// Send email here.
		if err := SendEmail(td, ctx, email, content); err != nil {
			errChan <- e.Wrap(err, "failed to send verification email")
		}
	}()

	go func() {
		defer wg.Done() // Notify the WaitGroup that this goroutine is done
		// Add link to db
		// TODO: If there is any other active code that hasn't expired invalidate all before creating another
		if err := store.CreateEmailVerificationRequest(ctx, sqlc.CreateEmailVerificationRequestParams{
			UserID:    userId,
			Email:     email,
			Token:     code,
			ExpiresAt: time.Now().Add(time.Minute * 15),
		}); err != nil {
			errChan <- e.Wrap(err, "failed to create email verification request")
		}
	}()

	// Wait for both goroutines to complete
	wg.Wait()
	close(errChan) // Close the channel so that the range loop can finish

	// Collect errors from the error channel
	for err := range errChan {
		if err != nil {
			return err // Return the first error encountered
		}
	}

	return nil

}

func VerifyEmail(ctx context.Context, store *sqlc.Store, code string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10) // Adjust the timeout as needed
	defer cancel()

	if len(code) != length {
		return errors.New("invalid token")
	}

	linkData, err := store.GetEmailVerificationRequestByToken(context.Background(), code)
	if err != nil {
		return err
	}

	if condition := linkData.ExpiresAt.Before(time.Now()); condition {
		return fmt.Errorf("token expired")
	}

	if condition := linkData.IsVerified.Bool; condition {
		return fmt.Errorf("token already verified")
	}

	// Update token to used
	usr, err := store.UpdateEmailVerificationRequest(context.Background(), sqlc.UpdateEmailVerificationRequestParams{
		Token:      code,
		IsVerified: sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		return err
	}

	// Verify user in "users" db
	_, err = store.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:              usr.UserID,
		IsEmailVerified: sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("error updating profile %s", err)
	}

	return nil

}

/*func SendVerificationEmailOnRegister(uid uuid.UUID, email string, username string, config utils.Config, store *sqlc.Store, ctx *gin.Context, l *zap.Logger) (string, error) {

	verifyCode, err := utils.GenerateUniqueToken(emailTokenLen)
	if err != nil {
		return "", err

	}

	link := config.AppUrl + "/verify?token=" + verifyCode
	msg := "Hello " + username + ", please verify your email address" + "with this link.\n" + link
	fmt.Println(msg)
	fmt.Println("Sent to ", email)

	// send email here.

	//   "is_verified" boolean DEFAULT false,
	//   "created_at" timestamptz DEFAULT (now()),

	// Add link to db
	err = store.CreateEmailVerificationRequest(context.Background(), sqlc.CreateEmailVerificationRequestParams{
		UserID:     uid,
		Email:      email,
		Token:      verifyCode, //link,
		IsVerified: sql.NullBool{Valid: true, Bool: false},
		ExpiresAt:  time.Now().Add(time.Minute * 15),
	})
	if err != nil {
		l.Error("Error creating email verification request", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "an unexpected error"})
		ctx.Abort()

		return "", err

	}

	return link, nil

}
*/

func checkAccountStatusForEmail(usr sqlc.Authentication) error {
	if usr.IsSuspended.Bool {
		return errors.New("account suspended")
	}

	if usr.IsDeleted.Bool {
		return errors.New("account deleted")
	}

	if usr.IsEmailVerified.Bool {
		return errors.New("email already verified")
	}
	return nil
}
