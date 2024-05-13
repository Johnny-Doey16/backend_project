package gapi

import (
	"context"
	"database/sql"
	"sync"

	"time"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	customError "github.com/steve-mir/diivix_backend/errors"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/pb"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/services"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) RotateToken(ctx context.Context, req *pb.RotateTokenRequest) (*pb.RotateTokenResponse, error) {
	// TODO: Verify that the access token is attached to the refresh token to avoid using an access token for a different refresh token (session)
	tokenMaker, err := token.NewPasetoMaker(s.config.RefreshTokenSymmetricKey)
	if err != nil {
		return nil, customError.New(codes.Internal, "cannot create token maker: %v", err)
	}

	payload, err := services.VerifyToken(tokenMaker, req.GetRefreshToken())
	if err != nil {
		return nil, customError.New(codes.Unauthenticated, "token verification failed: %v", err)
	}

	err = checkUserStatus(ctx, s.store, payload.UserId)
	if err != nil {
		return nil, err
	}

	session, err := s.store.GetSessionAndUserByRefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		if err == sql.ErrNoRows {
			if blockErr := blockUser(ctx, s.store, payload.UserId); blockErr != nil {
				return nil, blockErr
			}
			return nil, customError.New(codes.Internal, "suspicious activity detected")
		}
		return nil, customError.New(codes.Internal, "failed to get session: %v", err)
	}

	if !session.BlockedAt.Time.IsZero() && (session.BlockedAt.Time.Before(time.Now()) || session.InvalidatedAt.Time.Before(time.Now())) {
		return nil, customError.New(codes.PermissionDenied, "session blocked")
	}

	authToken, err := services.NewTokenService(s.config).
		RotateToken(session.Email, session.Username.String, session.Phone.String, true, session.IsEmailVerified.Bool, payload.UserId,
			int8(session.RoleID.Int32), session.ID, payload.IP, payload.UserAgent, s.config, s.store)

	if err != nil {
		return nil, customError.New(codes.Internal, "could not rotate token: %v", err)
	}

	return &pb.RotateTokenResponse{
		AuthInfo: &pb.AuthTokenInfo{
			AccessToken:           authToken.AccessToken,
			RefreshToken:          authToken.RefreshToken,
			AccessTokenExpiresAt:  timestamppb.New(authToken.AccessTokenExpiresAt),
			RefreshTokenExpiresAt: timestamppb.New(authToken.RefreshTokenExpiresAt),
		},
	}, nil
}

func checkUserStatus(ctx context.Context, store *sqlc.Store, uid uuid.UUID) error {
	user, err := store.GetUserByID(ctx, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return customError.New(codes.NotFound, "user not found")
		}
		return customError.New(codes.Internal, "failed to get user: %v", err)
	}

	if user.IsDeleted.Bool {
		return customError.New(codes.NotFound, "account not found")
	}

	if user.IsSuspended.Bool {
		return customError.New(codes.PermissionDenied, "account suspended")
	}
	return nil
}

func blockUser(ctx context.Context, store *sqlc.Store, userId uuid.UUID) error {
	// Create a channel to receive error results.
	errCh := make(chan error, 2) // Buffer size of 2 since we have two concurrent operations.
	var wg sync.WaitGroup

	// Start goroutine to block the user.
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := store.BlockUser(ctx, userId)
		errCh <- err
	}()

	// Start goroutine to block all user sessions.
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := store.BlockAllUserSession(ctx, userId)
		errCh <- err
	}()

	// Wait for the goroutines to finish.
	wg.Wait()
	close(errCh)

	// Check for errors.
	for err := range errCh {
		if err != nil {
			return customError.New(codes.Internal, "block operation error: %v", err)
		}
	}

	return nil
}
