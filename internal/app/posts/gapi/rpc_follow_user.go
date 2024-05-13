package gapi

import (
	"context"
	"fmt"

	// serv "github.com/steve-mir/diivix_backend/internal/app/authentication/services"
	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/constants"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *SocialMediaServer) FollowUser(ctx context.Context, req *pb.FollowUserRequest) (*pb.FollowResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	// followerID, err := services.StrToUUID("e7679d8b-0eac-4ea2-93cd-0018ab995922")
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "parsing uid: %s", err)
	// }
	followingID, err := services.StrToUUID(req.GetFollowingId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "parsing uid: %s", err)
	}

	err = services.Follow(ctx, s.db, s.store, claims.UserId, followingID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "following err: %s", err)
	}

	services.SendNotification(s.taskDistributor, ctx, claims.UserId, uuid.Nil, []uuid.UUID{followingID}, constants.NotificationFollow, "Follow", "", fmt.Sprintf("%s just followed you", claims.Username), "")

	// Write sql query to update user_profiles adding following_count and followers_count
	return &pb.FollowResponse{
		Success: true,
	}, nil
}
