package gapi

import (
	"context"
	"database/sql"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *SocialMediaServer) BlockUser(ctx context.Context, req *pb.BlockUserRequest) (*pb.BlockUserResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	blockedUid, err := services.StrToUUID(req.GetBlockedUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "parsing uid: %s", err)
	}

	blockedUser, err := s.store.BlockUserSM(ctx, sqlc.BlockUserSMParams{
		Reason:         sql.NullString{String: req.GetReason(), Valid: req.Reason != nil},
		BlockingUserID: claims.UserId,
		BlockedUserID:  blockedUid,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error blocking user: %s", err)
	}
	// Convert the sqlc BlockedUser type to the gRPC BlockedUser type.
	return &pb.BlockUserResponse{
		BlockedUser: convertToGrpcBlockedUser(blockedUser),
	}, nil
}

// RPC to get a list of blocked users by the blocker.
func (s *SocialMediaServer) GetBlockedUsers(ctx context.Context, req *pb.GetBlockedUsersRequest) (*pb.GetBlockedUsersResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	// Call the sqlc generated function to retrieve blocked users.
	blockedUsers, err := s.store.GetBlockedUsersByBlocker(ctx, claims.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting blocked users: %s", err)
	}
	// Convert the list of sqlc Blocked
	grpcBlockedUsers := make([]*pb.BlockedUser, len(blockedUsers))
	for i, bu := range blockedUsers {
		grpcBlockedUsers[i] = convertToGrpcBlockedUser(bu)
	}
	return &pb.GetBlockedUsersResponse{BlockedUsers: grpcBlockedUsers}, nil
}

// RPC to unblock a user.
func (s *SocialMediaServer) UnblockUser(ctx context.Context, req *pb.UnblockUserRequest) (*pb.UnblockUserResponse, error) {
	// Call the sqlc generated function to unblock the user.
	err := s.store.UnblockUser(ctx, req.Id)
	if err != nil {
		// Handle error.
		return nil, err
	}
	return &pb.UnblockUserResponse{
		Success: true,
	}, nil
}

// Helper function to convert sqlc BlockedUser type to gRPC BlockedUser type.
func convertToGrpcBlockedUser(bu sqlc.BlockedUser) *pb.BlockedUser {
	return &pb.BlockedUser{
		Id:             bu.ID,
		BlockingUserId: bu.BlockingUserID.String(),
		BlockedUserId:  bu.BlockedUserID.String(),
		Reason:         bu.Reason.String,
		CreatedAt:      timestamppb.New(bu.CreatedAt.Time),
	}
}
