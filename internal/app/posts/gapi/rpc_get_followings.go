package gapi

import (
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *SocialMediaServer) GetFollowing(req *pb.GetFollowingRequest, stream pb.SocialMedia_GetFollowingServer) error {
	// Parse the string representation of UUID
	parsedUUID, err := services.StrToUUID("fc5d5474-cc0a-4b98-831e-b1bb1b886580") // Assuming you get the UserId from the request
	if err != nil {
		return status.Errorf(codes.Internal, "parsing uid: %s", err)
	}

	// Retrieve the list of followers from the database
	followings, err := s.store.GetFollowing(stream.Context(), parsedUUID)
	if err != nil {
		return status.Errorf(codes.Internal, "error fetching followers: %s", err)
	}

	// Stream each follower profile back to the client
	for _, following := range followings {
		followerProfile := &pb.PostUser{
			// Assuming FollowerProfile has the fields Username, FirstName, and ImageUrl.
			// You'll need to replace these with the actual field names and sources from your follower struct.
			Uid:        following.ID.String(),
			Username:   following.Username.String,
			FirstName:  following.FirstName.String,
			ImageUrl:   following.ImageUrl.String,
			IsVerified: following.IsVerified.Bool,
			CreatedAt:  timestamppb.New(following.CreatedAt.Time),
		}

		if err := stream.Send(followerProfile); err != nil {
			return status.Errorf(codes.Internal, "error sending follower profile: %s", err)
		}
	}

	return nil
}
