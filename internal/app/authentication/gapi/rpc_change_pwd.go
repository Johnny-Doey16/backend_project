package gapi

import (
	"context"
	"errors"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/pb"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	if violations := validateChangePwdRequest(req); violations != nil {
		return nil, invalidArgumentErr(violations)
	}

	// Get current user
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	// TODO: Check if the current session (claims.SessionID) is active

	user, err := s.store.GetUserByID(ctx, claims.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't verify old password %s", err)
	}

	err = utils.CheckPassword(req.GetOldPassword(), user.PasswordHash)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "old password incorrect %s", err)
	}

	// Hash password
	pwdHash, err := utils.HashPassword(req.GetNewPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't verify new password %s", err)
	}

	// Check if new password is the same as the old password
	if req.GetOldPassword() == req.GetNewPassword() {
		return nil, status.Error(codes.InvalidArgument, "new password cannot be the same as the old password")
	}

	// Update password
	err = s.store.UpdateUserPassword(ctx, sqlc.UpdateUserPasswordParams{
		ID:           user.ID,
		PasswordHash: pwdHash,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't update password password %s", err)
	}

	// TODO: Close all the users session after changing password

	return &pb.ChangePasswordResponse{
		Success: true,
		Message: "password changed successfully",
	}, nil
}

func validateChangePwdRequest(req *pb.ChangePasswordRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if !utils.ValidatePassword(req.GetNewPassword()) {
		violations = append(violations, fieldViolation("new_password", errors.New("invalid password format")))
	}

	if req.GetNewPassword() != req.GetConfirmPassword() {
		violations = append(violations, fieldViolation("confirm_password", errors.New("both passwords must be same")))
	}

	return violations
}
