package gapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/cache"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *SocialMediaServer) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.PostStreamResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	postId, err := services.StrToUUID(req.PostId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "parsing uid: %s", err)
	}

	data, err := s.store.GetPostById(ctx, postId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error fetching post: %s", err)
	}

	return &pb.PostStreamResponse{
		Data: &pb.PostStreamResponse_Post{
			Post: &pb.Post{
				PostId:           data.PostID.String(),
				UserId:           data.UserID.String(),
				Content:          data.Content,
				Timestamp:        timestamppb.New(data.PostCreatedAt.Time),
				Username:         data.Username.String,
				FirstName:        "",
				ProfileImage:     data.ImageUrl.String,
				ImageUrls:        convertImgsToStr(&data),
				IsVerified:       data.IsVerified.Bool,
				FollowingsCounts: int64(data.FollowingCount.Int32),
				FollowersCounts:  int64(data.FollowersCount.Int32),
				Metrics: &pb.PostMetrics{
					Likes:    getPostLikeCounts(ctx, s.redisCache, PostLikesKey, &data), //data.Likes.Int32,
					Comments: data.Comments.Int32,
					Views:    data.Views.Int32,
					Reposts:  data.Reposts.Int32,
					Liked:    userLikedPost(ctx, s.redisCache, PostLikesKey, &data, claims.UserId),
				},
			},
		},
	}, nil
}

func convertImgsToStr(post *sqlc.GetPostByIdRow) []string {
	var imageUrls []string

	err := json.Unmarshal(post.ImageUrls, &imageUrls)
	if err != nil {

		log.Printf("cannot unmarshal data: %s", err)
	}
	return imageUrls
}

func userLikedPost(ctx context.Context, redisCache cache.Cache, postLikesKey string, post *sqlc.GetPostByIdRow, uid uuid.UUID) bool {
	// likeCountStr, _ := redisCache.GetKey(ctx, fmt.Sprintf("%s:%s", postLikesKey, post.ID))

	userLikesKey := fmt.Sprintf("%s:%s:users", postLikesKey, post.PostID)

	liked, _ := redisCache.SIsMember(ctx, userLikesKey, uid.String())
	return liked
}

func getPostLikeCounts(ctx context.Context, redisCache cache.Cache, postLikesKey string, post *sqlc.GetPostByIdRow) int32 {
	likeCountStr, _ := redisCache.GetKey(ctx, fmt.Sprintf("%s:%s", postLikesKey, post.PostID))
	return convertStrToInt32(likeCountStr)
}
