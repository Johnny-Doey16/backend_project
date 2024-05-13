package gapi

import (
	"context"

	"github.com/steve-mir/diivix_backend/internal/app/authentication/pb"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) RequestPasswordReset(ctx context.Context, req *pb.RequestPasswordResetRequest) (*pb.RequestPasswordResetResponse, error) {
	err := services.RequestPwdReset(ctx, s.store, s.taskDistributor, req.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error requesting password reset %v", err)
	}
	return &pb.RequestPasswordResetResponse{
		Message: services.ResetMsg,
	}, nil
}

func (s *Server) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	err := services.ResetPassword(ctx, s.db, s.store, req.GetToken(), req.GetNewPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not change password %v", err)
	}

	return &pb.ResetPasswordResponse{
		Message: "Password changed successfully",
	}, nil
}
