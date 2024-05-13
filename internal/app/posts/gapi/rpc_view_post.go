package gapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/cache"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *SocialMediaServer) ViewPost(ctx context.Context, req *pb.ViewPostRequest) (*pb.ViewPostResponse, error) {
	uid, _ := services.StrToUUID("ad7a61cd-14c5-4bbe-a3fe-0abdb585898a")

	postID, err := services.StrToUUID(req.GetPostId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Post ID incorrect: %v", err.Error())
	}

	// Create transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transaction error: %v", err.Error())
	}
	defer tx.Rollback()
	qtx := s.store.WithTx(tx)

	// Like in redis
	userViewsKey := fmt.Sprintf("post_views:%s:users", req.GetPostId())

	// Check if the user has already liked the post
	liked, err := s.redisCache.SIsMember(ctx, userViewsKey, uid.String())
	if err != nil {
		log.Printf("Error checking if user already liked the post: %v", err)
		return nil, err
	}

	var viewCount int64
	if !liked {
		// Mark post as viewed if the user has not viewed it
		viewCount, err = viewPost(ctx, qtx, tx, s.redisCache, userViewsKey, uid, postID)
		if err != nil {
			return nil, err
		}

	} else {
		log.Println("User has already viewed the post")
	}

	if err != nil {
		log.Printf("Error updating the like status for post: %v", err)
		return nil, err
	}

	err = broadcastViewEvent(ctx, s.redisCache, req.GetPostId(), uid)
	if err != nil {
		return nil, err
	}

	return &pb.ViewPostResponse{
		PostId:    req.PostId,
		ViewCount: viewCount,
	}, nil
}

func viewPost(ctx context.Context, qtx *sqlc.Queries, tx *sql.Tx, redisCache cache.Cache,
	userViewKey string, uid, postID uuid.UUID) (int64, error) {
	var viewCount int64
	err := redisCache.SAdd(ctx, userViewKey, uid.String())
	if err != nil {
		return 0, status.Errorf(codes.Internal, "Failed to add view to Redis: %v", err.Error())
	}

	err = services.ViewPostInDbConcurrently(tx, qtx, ctx, postID, uid)
	if err != nil {
		return 0, status.Errorf(codes.Internal, "Failed to view post in DB: %v", err.Error())
	}

	viewCount, err = redisCache.Incr(ctx, fmt.Sprintf("post_views:%s", postID.String()))
	if err != nil {
		return 0, status.Errorf(codes.Internal, "Failed to increment view count in Redis: %v", err.Error())
	}
	return viewCount, nil
}

func broadcastViewEvent(ctx context.Context, redisCache cache.Cache,
	postId string, uid uuid.UUID) error {
	likeUpdate := services.ViewUpdate{
		PostID: postId,
		UserID: uid.String(),
	}
	likeUpdateBytes, err := json.Marshal(likeUpdate)
	if err != nil {
		log.Printf("Error marshalling like update: %v", err)
		return err
	}

	// Publish the like update to the Redis channel
	err = redisCache.Publish(ctx, ViewChannelKey, likeUpdateBytes).Err()
	if err != nil {
		log.Printf("Error publishing like update: %v", err)
		return err
	}
	return nil
}
