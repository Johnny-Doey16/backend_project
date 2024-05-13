package gapi

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/cache"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *SocialMediaServer) GetPostsByUserId(req *pb.GetPostsByUserIdRequest, stream pb.SocialMedia_GetPostsByUserIdServer) error {
	// Parse the string representation of UUID
	parsedUUID, err := services.StrToUUID(req.UserId) // Assuming you get the UserId from the request
	if err != nil {
		return status.Errorf(codes.Internal, "parsing uid: %s", err)
	}

	offset := (req.GetPageNumber() - 1) * req.GetPageSize()
	// Note: Ensure that the page number and page size are positive.
	if req.PageNumber < 1 {
		req.PageNumber = 1
	}
	if req.GetPageSize() < 1 {
		req.PageSize = 10 // Default page size to 10 if not specified or negative
	}

	// Retrieve the list of followers from the database
	posts, err := s.store.GetPostByUserID(stream.Context(), sqlc.GetPostByUserIDParams{
		Limit:  req.GetPageSize(),
		Offset: int32(offset),
		UserID: parsedUUID,
	})
	if err != nil {
		return status.Errorf(codes.Internal, "error fetching followers: %s", err)
	}

	// Stream each follower profile back to the client
	for _, post := range posts {
		followerProfile := convertToPbPost(s.redisCache, &post, parsedUUID, PostLikesKey)

		if err := stream.Send(followerProfile); err != nil {
			return status.Errorf(codes.Internal, "error sending follower profile: %s", err)
		}
	}

	return nil
}

func convertToPbPost(redisCache cache.Cache, post *sqlc.GetPostByUserIDRow, uid uuid.UUID, postLikesKey string) *pb.Post {
	var imageUrls []string

	err := json.Unmarshal(post.PostImageUrls, &imageUrls)
	if err != nil {
		// Handle the error
		log.Printf("cannot unmarshal data: %s", err)
	}

	/*likeCountStr, _ := redisCache.GetKey(ctx, fmt.Sprintf("%s:%s", postLikesKey, post.ID))
	// Like in redis
	userLikesKey := fmt.Sprintf("%s:%s:users", postLikesKey, post.ID)

	// Check if the user has already liked the post
	liked, _ := redisCache.SIsMember(ctx, userLikesKey, uid.String())
	*/

	// Convert the sqlc.GetPostsRecommendationRow to the protobuf Post message
	return &pb.Post{
		UserId:       post.UserID.String(),
		Username:     post.Username.String,
		ProfileImage: post.UserImageUrl.String,
		Content:      post.Content,
		PostId:       post.ID.String(),
		ImageUrls:    imageUrls,
		Metrics: &pb.PostMetrics{
			Likes:    post.Likes.Int32,
			Views:    post.Views.Int32,
			Comments: post.Comments.Int32,
			Reposts:  post.Reposts.Int32,
			Liked:    post.PostLiked, //liked,
		},
		Timestamp:        timestamppb.New(post.CreatedAt.Time),
		IsVerified:       post.IsVerified.Bool,
		FollowingsCounts: int64(post.FollowingCount.Int32),
		FollowersCounts:  int64(post.FollowersCount.Int32),
	}
}

func convertStrToInt32(val string) int32 {
	if val == "" {
		return 0
	}
	// Convert string to int32
	int32Value, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}
	// Use the int32 value
	return int32(int32Value)
}
