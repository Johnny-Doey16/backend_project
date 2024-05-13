package gapi

import (
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *SocialMediaServer) GetFollowers(req *pb.GetFollowersRequest, stream pb.SocialMedia_GetFollowersServer) error {
	// Parse the string representation of UUID
	parsedUUID, err := services.StrToUUID("fc5d5474-cc0a-4b98-831e-b1bb1b886580") // Assuming you get the UserId from the request
	if err != nil {
		return status.Errorf(codes.Internal, "parsing uid: %s", err)
	}

	// Retrieve the list of followers from the database
	followers, err := s.store.GetFollowers(stream.Context(), parsedUUID)
	if err != nil {
		return status.Errorf(codes.Internal, "error fetching followers: %s", err)
	}

	// Stream each follower profile back to the client
	for _, follower := range followers {
		followerProfile := &pb.PostUser{
			// Assuming FollowerProfile has the fields Username, FirstName, and ImageUrl.
			// You'll need to replace these with the actual field names and sources from your follower struct.
			Uid:        follower.ID.String(),
			Username:   follower.Username.String,
			FirstName:  follower.FirstName.String,
			ImageUrl:   follower.ImageUrl.String,
			IsVerified: follower.IsVerified.Bool,
			CreatedAt:  timestamppb.New(follower.CreatedAt.Time),
		}

		if err := stream.Send(followerProfile); err != nil {
			return status.Errorf(codes.Internal, "error sending follower profile: %s", err)
		}
	}

	return nil
}
