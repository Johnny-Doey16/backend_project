package gapi

import (
	"context"
	"time"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *SocialMediaServer) SuggestUsersToFollow(ctx context.Context, req *pb.SuggestUsersToFollowRequest) (*pb.SuggestUsersToFollowResponse, error) {
	// claims, ok := ctx.Value("payloadKey").(*token.Payload)
	// if !ok {
	//     return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	// }

	userId, err := services.StrToUUID("558c60aa-977f-4d38-885b-e813656371ac") // ad7a61cd-14c5-4bbe-a3fe-0abdb585898a
	if err != nil {
		return nil, status.Errorf(codes.Internal, "parsing uid: %s", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*10) // Adjust the timeout as needed
	defer cancel()

	// Calculate the offset based on the page number and page size
	offset := (req.GetPageNumber() - 1) * req.GetPageSize()
	limit := req.GetPageSize()

	// Calculate the offset based on the page number and page size.
	// Note: Ensure that the page number and page size are positive.
	if req.PageNumber < 1 {
		req.PageNumber = 1
	}
	if limit < 1 {
		limit = 10 // Default page size to 10 if not specified or negative
	}

	users, err := s.store.SuggestUsersToFollowPaginated(ctx, sqlc.SuggestUsersToFollowPaginatedParams{
		ID:     userId,
		Limit:  limit,
		Offset: int32(offset),
	})
	// users, err := s.store.SuggestUsersToFollowPaginated(ctx, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error fetching suggested users %v", err.Error())
	}

	// Create response from the fetched users.
	resp := &pb.SuggestUsersToFollowResponse{
		Users:   []*pb.SuggestedUsers{}, // Assuming SuggestedUser is the correct message name
		HasMore: true,
	}
	for _, user := range users {
		resp.Users = append(resp.Users, &pb.SuggestedUsers{
			Id:             user.ID.String(),     //user.UserID.String(), //.ID.String(),
			Username:       user.Username.String, // Assuming the Username field should map to FirstName
			ImageUrl:       user.ImageUrl.String,
			FirstName:      user.FirstName.String,
			LastName:       user.LastName.String,
			IsVerified:     user.IsVerified.Bool,
			FollowingCount: int64(user.FollowingCount.Int32),
			FollowerCount:  int64(user.FollowersCount.Int32),
		})
	}

	// Set HasMore to false if the number of returned users is less than the page size
	if len(users) < int(limit) {
		resp.HasMore = false
	}

	return resp, nil
}
