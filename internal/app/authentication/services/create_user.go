package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sqlc-dev/pqtype"
	"github.com/steve-mir/diivix_backend/cache"
	"github.com/steve-mir/diivix_backend/constants"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/pb"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/utils"
	ut "github.com/steve-mir/diivix_backend/utils"
	"github.com/steve-mir/diivix_backend/worker"
)

type UserResult struct {
	User sqlc.Authentication
	Err  error
}

type AccessTokenResult struct {
	AccessToken string
	Payload     *token.Payload
	Err         error
}

func CheckUserExists(ctx context.Context, qtx *sqlc.Queries, email, username string) error {
	// Check db if email exists
	if err := checkEmailExistsError(ctx, qtx, email); err != nil {
		return err
	}

	// Check db if username exists
	if err := CheckIfUsernameExists(ctx, qtx, username); err != nil {
		return err
	}
	return nil
}

func PrepareUserData(pwd string) (string, uuid.UUID, error) {
	hashedPwd, err := utils.HashPassword(pwd)
	if err != nil {
		return "", uuid.UUID{}, errors.New("error processing data")
	}

	// Generate UUID in advance
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", uuid.UUID{}, errors.New("an unexpected error occurred")
	}

	return hashedPwd, uid, nil
}

func CreateUserConcurrent(ctx context.Context, qtx *sqlc.Queries, tx *sql.Tx, uid uuid.UUID, email, username, pwd string) (sqlc.Authentication, error) {
	params := sqlc.CreateUserParams{
		ID:           uid,
		Email:        email,
		Username:     sql.NullString{String: username, Valid: true},
		PasswordHash: pwd,
		UserType:     constants.RegularUsersUser,
	}

	var wg sync.WaitGroup
	createUserChan := make(chan UserResult, 1)

	// Start the post creation goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		user, err := qtx.CreateUser(ctx, params)
		createUserChan <- UserResult{User: user, Err: err}
	}()

	// Wait for the post creation to complete before starting other operations
	wg.Wait()
	userData := <-createUserChan
	if userData.Err != nil {
		tx.Rollback()
		return sqlc.Authentication{}, errors.New("error creating post " + userData.Err.Error())
	}

	return userData.User, nil
}

func RunConcurrentUserCreationTasks(ctx context.Context, qtx *sqlc.Queries, tx *sql.Tx, config ut.Config, td worker.TaskDistributor,
	req *pb.CreateUserRequest, uid uuid.UUID, clientIP pqtype.Inet, agent string) (string, time.Time, error) {
	var wg sync.WaitGroup
	wg.Add(6)

	createAccessTokenChan := make(chan AccessTokenResult, 1)
	createProfileChan := make(chan error, 1)
	createRoleChan := make(chan error, 1)
	sendEmailChan := make(chan error, 1)
	createUserNames := make(chan error, 1)
	createUserAccountChan := make(chan error, 1)

	go func() {
		defer wg.Done()
		tokenService := NewTokenService(config)
		accessToken, accessPayload, err := tokenService.CreateAccessToken(req.GetEmail(), req.GetUsername(), "", true, false, uid, constants.RegularUsers, clientIP, agent)
		createAccessTokenChan <- AccessTokenResult{AccessToken: accessToken, Payload: accessPayload, Err: err}
	}()

	go func() {
		defer wg.Done()
		// TODO: Fix to add names
		profileErr := qtx.CreateUserNames(ctx, sqlc.CreateUserNamesParams{
			UserID: uid,
		})
		createUserNames <- profileErr
	}()

	go func() {
		defer wg.Done()
		profileErr := qtx.CreateEntityProfile(ctx, sqlc.CreateEntityProfileParams{
			UserID:     uid,
			EntityType: constants.RegularUsersUser,
		})
		createProfileChan <- profileErr
	}()

	go func() {
		defer wg.Done()
		_, roleErr := qtx.CreateUserRole(ctx, sqlc.CreateUserRoleParams{
			UserID: uid,
			RoleID: constants.RegularUsers,
		})
		createRoleChan <- roleErr
	}()

	go func() {
		defer wg.Done()
		err := qtx.CreateUserAccountDetails(ctx, sqlc.CreateUserAccountDetailsParams{
			UserID: uid,
		})
		createUserAccountChan <- err
	}()

	go func() {
		defer wg.Done()
		err := SendVerificationEmail(qtx, ctx, td, uid, req.GetEmail())
		sendEmailChan <- err
	}()

	wg.Wait()
	close(createAccessTokenChan)
	close(createProfileChan)
	close(createRoleChan)
	close(createUserNames)
	close(createUserAccountChan)
	close(sendEmailChan)

	claims := <-createAccessTokenChan
	if claims.Err != nil {
		tx.Rollback()
		return "", time.Time{}, errors.New("unknown error")
	}

	if err := <-createUserNames; err != nil {
		tx.Rollback()
		return "", time.Time{}, fmt.Errorf("an unknown error occurred creating users %v", err)
	}

	if err := <-createProfileChan; err != nil {
		tx.Rollback()
		return "", time.Time{}, fmt.Errorf("an unknown error occurred creating profile %v", err)
	}

	if <-createRoleChan != nil {
		tx.Rollback()
		return "", time.Time{}, errors.New("error cannot proceed")
	}

	if err := <-createUserAccountChan; err != nil {
		tx.Rollback()
		return "", time.Time{}, errors.New("error creating user account details " + err.Error())
	}

	if err := <-sendEmailChan; err != nil {
		tx.Rollback()
		return "", time.Time{}, errors.New("unable to resend email " + err.Error())
	}

	return claims.AccessToken, claims.Payload.Expires, nil
}

func CreateNewUser(ctx context.Context, qtx *sqlc.Queries, uid uuid.UUID, email, username, pwd, userType string) (sqlc.Authentication, error) {
	params := sqlc.CreateUserParams{
		ID:           uid,
		Email:        email,
		Username:     sql.NullString{String: username, Valid: true},
		PasswordHash: pwd,
		UserType:     userType,
	}

	return qtx.CreateUser(ctx, params)
}

// ?----------------
func checkEmailExistsError(ctx context.Context, qtx *sqlc.Queries, email string) error {
	// Check duplicate emails
	user, err := qtx.GetUserByIdentifier(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		// An error occurred that isn't simply indicating no rows were found
		return err
	}

	if user.ID != uuid.Nil {
		// User exists, check if the account is marked as deleted
		if user.DeletedAt.Valid {
			// Check if the account is within the recovery period
			if time.Since(user.DeletedAt.Time) <= MaxAccountRecoveryDuration {
				// Account is within the recovery period and can be recovered
				return errors.New("account is deleted but can be recovered, please follow the account recovery process")
			} else {
				// Account is beyond the recovery period, append timestamp to the email to make it unique
				err = appendTimestampToEmail(ctx, qtx, user.Email, user.DeletedAt.Time)
				if err != nil {
					return fmt.Errorf("failed to update email for user with ID %s: %v", user.ID, err)
				}
			}
		} else {
			// Account exists and is not marked as deleted
			return errors.New("email already exists")
		}
	}
	return nil
}

// Placeholder store method to append a timestamp to the user's email
func appendTimestampToEmail(ctx context.Context, qtx *sqlc.Queries, email string, deletedAt time.Time) error {
	// Implement the logic to append a timestamp to the user's email.
	// This will involve updating the user record in the database.
	// Be careful to ensure that the new email remains unique and valid.
	// For example, you might append something like "_deleted_1612385610" to the email.
	newEmail := addDeleteTimeToEmail(email, deletedAt)

	_, err := qtx.UpdateUser(ctx, sqlc.UpdateUserParams{
		Email: sql.NullString{String: newEmail, Valid: true},
	})
	if err != nil {
		return err
	}

	return nil
}

func addDeleteTimeToEmail(email string, deletedAt time.Time) string {
	timestamp := deletedAt.Unix() // Convert time to Unix timestamp
	modifiedEmail := fmt.Sprintf("%s_deleted_%d", email, timestamp)
	return modifiedEmail
}

func FetchUserSuggestions(redisCache cache.Cache, ctx context.Context, query string) ([]string, error) {
	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Initialize variables for iteration
	var cursor uint64 = 0
	var usernames []string

	// If query string is empty, return an empty result
	if query == "" {
		return usernames, nil
	}

	// Iterate over keyspace
	for {
		// Scan for keys matching the pattern "user:<query>*", limit to 10 keys per iteration
		keys, nextCursor, err := rdb.Scan(ctx, cursor, "user:"+query+"*", 10).Result()
		if err != nil {
			return nil, err
		}

		// Append usernames to result
		for _, key := range keys {
			username := strings.TrimPrefix(key, "user:")
			usernames = append(usernames, username)
		}

		// Update cursor for next iteration
		cursor = nextCursor

		// Break if iteration is complete
		if cursor == 0 {
			break
		}
	}

	return usernames, nil
}

func StoreUser(redisCache cache.Cache, ctx context.Context, username string) {
	// Connect to Redis
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "", // no password set
	// 	DB:       0,  // use default DB
	// })

	// Store user in Redis
	if err := redisCache.SetKey(ctx, "user:"+username, "", 0); err != nil {
		log.Fatalf("error storing user in Redis: %v", err)
	}
	log.Printf("User stored in Redis: %s", username)
}
