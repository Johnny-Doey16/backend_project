package gapi

import (
	"context"

	"github.com/steve-mir/diivix_backend/internal/app/authentication/pb"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/services"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) ResendVerification(ctx context.Context, _ *pb.ResendVerificationRequest) (*pb.ResendVerificationResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	err := services.ReSendVerificationEmail(s.store, ctx, s.taskDistributor, claims.UserId, claims.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to resend email "+err.Error())
	}

	return &pb.ResendVerificationResponse{
		Message: "email sent check your inbox",
	}, nil
}

// RPC to verify an email.
func (s *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	err := services.VerifyEmail(ctx, s.store, req.GetToken())
	if err != nil {
		return nil, status.Error(codes.Internal, "could not verify email, "+err.Error())
	}
	return &pb.VerifyEmailResponse{
		Message: "email verified successfully",
	}, nil
}
