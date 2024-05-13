package gapi

import (
	"context"
	"time"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *SocialMediaServer) FetchPostComment(ctx context.Context, req *pb.FetchCommentsRequest) (*pb.FetchCommentsResponse, error) {

	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	postID, err := services.StrToUUID(req.GetPostId())
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

	comments, err := s.store.RecommendComments(ctx, sqlc.RecommendCommentsParams{
		PostID:         postID,
		FollowerUserID: claims.UserId,
		Limit:          req.PageSize,
		Offset:         int32(offset),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error fetching suggested users %v", err.Error())
	}

	// Create response from the fetched users.
	resp := &pb.FetchCommentsResponse{
		Comments: []*pb.Comment{}, // Assuming SuggestedUser is the correct message name
		HasMore:  true,
	}
	for _, comment := range comments {
		resp.Comments = append(resp.Comments, &pb.Comment{
			CommentId:    comment.CommentID,
			UserId:       comment.UserID.String(),
			PostId:       req.PostId,
			CreatedAt:    timestamppb.New(comment.CommentCreatedAt.Time),
			UpdatedAt:    timestamppb.New(comment.CommentUpdatedAt.Time),
			Text:         comment.CommentText,
			ProfileImage: comment.ImageUrl.String,
			FirstName:    comment.FirstName.String,
			Username:     comment.CommenterUsername.String,
		})
	}

	// Set HasMore to false if the number of returned users is less than the page size
	if len(comments) < int(limit) {
		resp.HasMore = false
	}

	return resp, nil
}
