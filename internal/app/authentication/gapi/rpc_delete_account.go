package gapi

import (
	"context"
	"fmt"

	"github.com/steve-mir/diivix_backend/internal/app/authentication/pb"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/services"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	err := services.DeleteAccountRequest(ctx, req.GetPassword(), s.store, claims.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteAccountResponse{
		Message: fmt.Sprintf("%s deleted successfully", claims.Email),
	}, nil
}

// Initiates the process to recover a deleted account.
func (s *Server) RequestAccountRecovery(ctx context.Context, req *pb.RecoveryRequest) (*pb.RecoveryResponse, error) {
	err := services.AccRecoveryRequest(ctx, s.store, s.taskDistributor, req.GetEmail())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.RecoveryResponse{
		Message: "an email to recover you account has been sent",
		Success: true,
	}, nil
}

// Completes the account recovery process.
func (s *Server) CompleteAccountRecovery(ctx context.Context, req *pb.CompleteRecoveryRequest) (*pb.CompleteRecoveryResponse, error) {
	err := services.AccountRecovery(ctx, s.db, s.store, req.GetRecoveryToken())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.CompleteRecoveryResponse{
		Success: true,
		Message: "account successfully recovered",
	}, nil
}
