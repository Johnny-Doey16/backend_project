package gapi

import (
	"context"

	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *SocialMediaServer) UnFollowUser(ctx context.Context, req *pb.UnFollowUserRequest) (*pb.FollowResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	followingID, err := services.StrToUUID(req.GetFollowingId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "parsing uid: %s", err)
	}

	err = services.UnFollow(ctx, s.db, claims.UserId, followingID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "following err: %s", err)
	}

	// Write sql query to update user_profiles adding following_count and followers_count
	return &pb.FollowResponse{
		Success: true,
	}, nil
}
