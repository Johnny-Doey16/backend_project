package gapi

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/pb"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/services"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// RPC to register a new user for MFA.
func (s *Server) RegisterMFA(ctx context.Context, req *pb.RegisterMFARequest) (*pb.RegisterMFAResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	mfaSecret, qrCode, url, err := services.RegisterMFA(ctx, s.config, s.store, claims.UserId, s.config.AppName, req.GetPassword())
	if err != nil {
		return nil, err
	}

	// // set mfa enabled in auth to true
	// _, err = s.store.UpdateUser(ctx, sqlc.UpdateUserParams{
	// 	ID:           claims.UserId,
	// 	IsMfaEnabled: sql.NullBool{Bool: true, Valid: true},
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// Return mfaSecret, recoveryCodes and qrCode byte
	return &pb.RegisterMFAResponse{
		Secret: mfaSecret,
		QrCode: qrCode,
		Url:    &url,
		// RecoveryCodes: recoveryCodes,
	}, nil
}

// RPC to verify a TOTP code during mfa setup
func (s *Server) VerifyMFAWorks(ctx context.Context, req *pb.VerifyMFAWorksRequest) (*pb.VerifyMFAWorksResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	recoveryCode, err := services.ValidateMFAWorks(ctx, s.config, s.db, s.store, claims.UserId, req.GetToken(), req.GetSecret())
	if err != nil {
		fmt.Printf("Error %v", err)
		return nil, err
	}

	// set mfa enabled in auth to true
	_, err = s.store.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:           claims.UserId,
		IsMfaEnabled: sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	//  rotate access token with mfa passed
	// authToken, err := services.NewTokenService(s.config).
	// 	RotateToken(claims.Email, claims.Username, claims.Phone, true, claims.IsEmailVerified, claims.UserId,
	// 		claims.Role, claims.SessionID, claims.IP, claims.UserAgent, s.config, s.store)

	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "could not rotate token: %v", err)
	// }

	return &pb.VerifyMFAWorksResponse{
		RecoveryCodes: recoveryCode,
	}, nil
}

// RPC to verify a TOTP code during sign-in.
func (s *Server) VerifyMFA(ctx context.Context, req *pb.VerifyMFARequest) (*pb.VerifyMFAResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	err := services.ValidateMFA(ctx, s.config, s.store, claims.UserId, req.GetToken())
	if err != nil {
		fmt.Printf("Error %v", err)
		return nil, err
	}

	// ! rotate access token with mfa passed
	authToken, err := services.NewTokenService(s.config).
		RotateToken(claims.Email, claims.Username, claims.Phone, true, claims.IsEmailVerified, claims.UserId,
			claims.Role, claims.SessionID, claims.IP, claims.UserAgent, s.config, s.store)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not rotate token: %v", err)
	}

	return &pb.VerifyMFAResponse{
		Success: true,
		AuthInfo: &pb.AuthTokenInfo{
			AccessToken:           authToken.AccessToken,
			RefreshToken:          authToken.RefreshToken,
			AccessTokenExpiresAt:  timestamppb.New(authToken.AccessTokenExpiresAt),
			RefreshTokenExpiresAt: timestamppb.New(authToken.RefreshTokenExpiresAt),
		},
	}, nil
}

// RPC to generate an OTP for MFA setup.
func (s *Server) ByPassMFA(ctx context.Context, req *pb.ByPassOtpRequest) (*pb.ByPassOtpResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	err := services.BypassMFA(ctx, s.config, s.store, claims.UserId, req.GetRecoveryCode())
	if err != nil {
		return nil, err
	}

	// ! rotate access token with mfa by passed
	authToken, err := services.NewTokenService(s.config).
		RotateToken(claims.Email, claims.Username, claims.Phone, true, claims.IsEmailVerified, claims.UserId,
			claims.Role, claims.SessionID, claims.IP, claims.UserAgent, s.config, s.store)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not rotate token: %v", err)
	}

	// Generate tokens
	return &pb.ByPassOtpResponse{
		Success: true,
		AuthInfo: &pb.AuthTokenInfo{
			AccessToken:           authToken.AccessToken,
			RefreshToken:          authToken.RefreshToken,
			AccessTokenExpiresAt:  timestamppb.New(authToken.AccessTokenExpiresAt),
			RefreshTokenExpiresAt: timestamppb.New(authToken.RefreshTokenExpiresAt),
		},
	}, nil
}
