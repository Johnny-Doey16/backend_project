package gapi

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *SocialMediaServer) BookmarkPost(ctx context.Context, req *pb.BookmarkPostRequest) (*pb.BookmarkPostResponse, error) {
	// claims, ok := ctx.Value("payloadKey").(*token.Payload)
	// if !ok {
	// 	return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	// }
	uid, _ := services.StrToUUID("9cd93ffa-5835-414f-95c7-f5e94608c71c")

	// Convert the post ID from the request to UUID format
	postID, err := services.StrToUUID(req.GetPostId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid post ID format: %s", err.Error())
	}

	// Check post status
	postStatus, err := s.store.CheckPostStatus(ctx, postID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error checking post status: %s", err.Error())
	}

	// Check if the post is active
	if postStatus.PostStatus == "active" {

		// Attempt to add post to bookmarks
		err = s.store.CreateBookmark(ctx, sqlc.CreateBookmarkParams{
			UserID: uid,
			PostID: postID,
		})

		// Check for errors
		if err != nil {
			// Check if the error is due to a unique constraint violation
			if strings.Contains(err.Error(), "unique constraint") {
				return nil, status.Errorf(codes.AlreadyExists, "Post is already bookmarked")
			}
			// For other errors, return a generic internal server error response
			return nil, status.Errorf(codes.Internal, "Error creating bookmark: %s", err.Error())
		}

		// If no errors, return a success message
		return &pb.BookmarkPostResponse{
			Message: "Post bookmarked successfully",
			Success: true,
		}, nil
	}
	// If the post is not active, return an appropriate response
	return nil, status.Errorf(codes.FailedPrecondition, "Cannot bookmark a %s post", postStatus)

}

func (s *SocialMediaServer) GetBookmarkedPosts(ctx context.Context, req *pb.GetBookmarkedPostsRequest) (*pb.GetBookmarkedPostsResponse, error) {
	// claims, ok := ctx.Value("payloadKey").(*token.Payload)
	// if !ok {
	// 	return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	// }
	uid, _ := services.StrToUUID("9cd93ffa-5835-414f-95c7-f5e94608c71c")

	// Default page size if not specified or negative
	const defaultPageSize = 10

	pageNumber := req.GetPageNumber()
	if pageNumber < 1 {
		pageNumber = 1
	}

	pageSize := req.GetPageSize()
	if pageSize < 1 {
		pageSize = defaultPageSize
	}

	offset := (pageNumber - 1) * pageSize

	posts, err := s.store.GetUserBookmarks(ctx, sqlc.GetUserBookmarksParams{
		UserID: uid,
		Offset: int32(offset),
		Limit:  int32(pageSize),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error fetching bookmarks: %s", err.Error())
	}

	resp := &pb.GetBookmarkedPostsResponse{
		Posts:   []*pb.BookmarkedPost{},
		HasMore: len(posts) == int(pageSize),
	}

	for _, post := range posts {
		imageUrls := []string{}
		if err := json.Unmarshal(post.PostImageUrls, &imageUrls); err != nil {
			log.Printf("cannot unmarshal image URLs: %v", err)
			continue // Optionally, you can decide to handle this error differently
		}

		// Skip posts that are deleted or suspended
		if post.PostDeletedAt.Valid || post.PostSuspendedAt.Valid {
			continue
		}

		resp.Posts = append(resp.Posts, &pb.BookmarkedPost{
			Post: &pb.Post{
				Content:      post.Content,
				UserId:       post.AuthorUserID.String(),
				PostId:       post.PostID.String(),
				FirstName:    post.AuthorFirstName.String,
				Username:     post.AuthorUsername.String,
				ProfileImage: post.AuthorImageUrl.String,
				ImageUrls:    imageUrls,
				Timestamp:    timestamppb.New(post.PostCreatedAt.Time),
				Metrics: &pb.PostMetrics{
					Views:    post.Views.Int32,
					Likes:    post.Likes.Int32,
					Comments: post.Comments.Int32,
					Liked:    post.Liked,
				},
			},
			BookmarkedAt: timestamppb.New(post.BookmarkedAt.Time),
		})
	}

	return resp, nil
}

// RPC to delete a bookmarked post.
func (s *SocialMediaServer) DeleteBookmark(ctx context.Context, req *pb.DeleteBookmarkRequest) (*pb.DeleteBookmarkResponse, error) {
	// claims, ok := ctx.Value("payloadKey").(*token.Payload)
	// if !ok {
	// 	return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	// }
	uid, _ := services.StrToUUID("9cd93ffa-5835-414f-95c7-f5e94608c71c")

	// Start transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to begin transaction: %v", err)
	}
	defer tx.Rollback() // Rollback will be ignored if tx has been committed later on

	qtx := s.store.WithTx(tx)

	// Attempt to delete the bookmark
	bk, err := qtx.DeleteBookmark(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete bookmark: %v", err)
	}

	// Check if the bookmark belongs to the user
	if bk.UserID != uid {
		return nil, status.Error(codes.PermissionDenied, "Bookmark does not belong to the user")
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to commit deletion: %v", err)
	}

	return &pb.DeleteBookmarkResponse{
		Message: "Bookmark successfully deleted",
		Success: true,
	}, nil
}
