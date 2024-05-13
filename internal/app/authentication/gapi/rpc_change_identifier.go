package gapi

import (
	"context"

	"github.com/steve-mir/diivix_backend/internal/app/authentication/pb"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/services"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// * Research more on what happens to tokens with the old data after changing to the new credentials

func (s *Server) InitiateChangeEmail(ctx context.Context, req *pb.InitiateChangeEmailRequest) (*pb.InitiateChangeEmailResponse, error) {
	// Get current user
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	err := services.RequestChangeOfEmail(s.taskDistributor, ctx, s.store, claims.UserId, claims.Email, req.GetNewEmail())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.InitiateChangeEmailResponse{
		Message: "A code has been sent to your new email successfully",
	}, nil
}

// Confirm the change of email
func (s *Server) ConfirmChangeEmail(ctx context.Context, req *pb.ConfirmChangeEmailRequest) (*pb.ConfirmChangeEmailResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	accessToken, accessExpires, err := services.ChangeEmail(ctx, s.db, s.config, s.store, claims, req.VerificationCode)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.ConfirmChangeEmailResponse{
		Message:              "Changed successfully",
		AccessToken:          accessToken,
		AccessTokenExpiresAt: timestamppb.New(accessExpires),
	}, nil
}

// Initiate a change phone number request
func (s *Server) InitiateChangePhone(ctx context.Context, req *pb.InitiateChangePhoneRequest) (*pb.InitiateChangePhoneResponse, error) {

	// Get current user
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	err := services.RequestChangeOfPhone(ctx, s.store, claims.UserId, claims.Email, req.GetNewPhone())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.InitiateChangePhoneResponse{
		Message: "A code has been sent to your new phone successfully",
	}, nil
}

// Confirm the change of phone number
func (s *Server) ConfirmChangePhone(ctx context.Context, req *pb.ConfirmChangePhoneRequest) (*pb.ConfirmChangePhoneResponse, error) {

	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	accessToken, accessExpires, err := services.ChangePhone(ctx, s.db, s.config, s.store, claims, req.VerificationCode)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.ConfirmChangePhoneResponse{
		Message:              "Changed successfully",
		AccessToken:          accessToken,
		AccessTokenExpiresAt: timestamppb.New(accessExpires),
	}, nil
}

// Update username
func (s *Server) ChangeUsername(ctx context.Context, req *pb.ChangeUsernameRequest) (*pb.ChangeUsernameResponse, error) {

	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return &pb.ChangeUsernameResponse{Message: "Error.."}, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	accessToken, accessExpires, err := services.ChangeUsername(ctx, s.db, s.config, s.store, claims, req.GetNewUsername()) //ctx, s.db, s.store, claims.UserId, req.NewUsername)
	if err != nil {
		return &pb.ChangeUsernameResponse{Message: "Error: " + err.Error()}, err
	}

	return &pb.ChangeUsernameResponse{
		Message:              "username changed successfully",
		AccessToken:          accessToken,
		AccessTokenExpiresAt: timestamppb.New(accessExpires),
	}, nil
}
