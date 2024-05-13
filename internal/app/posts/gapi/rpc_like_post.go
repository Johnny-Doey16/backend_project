package gapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/cache"
	"github.com/steve-mir/diivix_backend/constants"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// abee399e-a61f-4846-b1f0-63998ce57d8f
func (s *SocialMediaServer) LikePost(ctx context.Context, req *pb.LikePostRequest) (*pb.LikePostResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}
	// userId, err := services.StrToUUID("558c60aa-977f-4d38-885b-e813656371ac") // ad7a61cd-14c5-4bbe-a3fe-0abdb585898a
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "parsing uid: %s", err)
	// }

	postID, postUserID, err := getIdsFromString(req.GetPostId(), req.GetPostUserId())
	if err != nil {
		return nil, err
	}

	// Create transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transaction error: %v", err.Error())
	}
	defer tx.Rollback()
	qtx := s.store.WithTx(tx)

	// Like in redis
	/*
		userLikesKey := fmt.Sprintf("post_likes:%s:users", req.GetPostId())
			// // Check if the user has already liked the post
			liked, err := s.redisCache.SIsMember(ctx, userLikesKey, userId)
			if err != nil {
				log.Printf("Error checking if user already liked the post: %v", err)
				return nil, err
			}

			var likeCount int64
			var postLiked bool
			if liked {
				log.Println("UnLiking post", postID, userId)
				// User already liked the post, remove the like
				likeCount, err = unLikePost(ctx, qtx, tx, s.redisCache, userLikesKey, userId, postID)
				if err != nil {
					return nil, err
				}

			} else {
				log.Println("Liking post", postID, userId)
				// User has not liked the post, add the like
				likeCount, postLiked, err = likePost(ctx, qtx, tx, s.redisCache, userLikesKey, userId, postID)
				if err != nil {
					return nil, err
				}
			}

			if err != nil {
				log.Printf("Error updating the like status for post: %v", err)
				return nil, err
			}
	*/

	likeCount, liked, err := likePost(ctx, qtx, tx, s.redisCache, claims.UserId, postID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to like post in DB: %v", err.Error())
	}

	// broadcast like event here
	// TODO: Broadcast liked too.
	if err := broadcastLikeEvent(ctx, s.redisCache, req.GetPostId(), claims.UserId); err != nil {
		return nil, err
	}

	// Send notification
	if postUserID != claims.UserId && liked {
		services.SendNotification(s.taskDistributor, ctx, claims.UserId, postID, []uuid.UUID{postUserID}, constants.NotificationPostLike, "Like", "", fmt.Sprintf("%s liked your post", claims.Username), "")
	}

	// Return a response indicating success
	return &pb.LikePostResponse{
		PostId:    req.GetPostId(),
		LikeCount: int64(likeCount), // Return the actual like count
		Like:      liked,
	}, nil
}

func likePost(ctx context.Context, qtx *sqlc.Queries, tx *sql.Tx, redisCache cache.Cache,
	uid, postID uuid.UUID) (int64, bool, error) {

	ls, liked, err := services.LikePostInDbConcurrently(tx, qtx, ctx, postID, uid)
	if err != nil {
		return 0, false, status.Errorf(codes.Internal, "Failed to like post in DB: %v", err.Error())
	}

	// Set the like count in Redis to match the database count
	likeCountKey := fmt.Sprintf("%s:%s", PostLikesKey, postID.String())
	err = redisCache.SetKey(ctx, likeCountKey, ls, time.Hour*24)
	if err != nil {
		return 0, false, status.Errorf(codes.Internal, "Failed to set like count in Redis: %v", err.Error())
	}

	return int64(ls), liked, nil
}

func broadcastLikeEvent(ctx context.Context, redisCache cache.Cache,
	postId string, uid uuid.UUID) error {
	likeUpdate := services.LikeUpdate{
		PostID: postId,
		UserID: uid.String(),
		// Additional fields as needed
	}
	likeUpdateBytes, err := json.Marshal(likeUpdate)
	if err != nil {
		log.Printf("Error marshalling like update: %v", err)
		return err
	}

	// Publish the like update to the Redis channel
	err = redisCache.Publish(ctx, LikeChannelKey, likeUpdateBytes).Err()
	log.Println("Broadcasting", LikeChannelKey)
	if err != nil {
		log.Printf("Error publishing like update: %v", err)
		return err
	}
	return nil
}

/*func unLikePost(ctx context.Context, qtx *sqlc.Queries, tx *sql.Tx, redisCache cache.Cache,
	userLikesKey string, uid, postID uuid.UUID) (int64, error) {

	err := redisCache.SRem(ctx, userLikesKey, uid.String())
	if err != nil {
		return 0, status.Errorf(codes.Internal, "Failed to remove like from Redis: %v", err.Error())
	}

	ls, err := services.UnLikePostInDbConcurrently(tx, qtx, ctx, postID, uid)
	if err != nil {
		return 0, status.Errorf(codes.Internal, "Failed to unlike post in DB: %v", err.Error())
	}

	// Set the like count in Redis to match the database count
	likeCountKey := fmt.Sprintf("%s:%s", PostLikesKey, postID.String())
	err = redisCache.SetKey(ctx, likeCountKey, ls, time.Hour*24)
	if err != nil {
		return 0, status.Errorf(codes.Internal, "Failed to set like count in Redis: %v", err.Error())
	}

	return int64(ls), nil
}
*/

// Concurrency Management: Locking and unlocking the subscribers map for every message can become a bottleneck and may lead to contention issues as the number of concurrent operations increases.

// In-Memory Data Stores: Using an in-memory data store like Redis with Pub/Sub capabilities can facilitate broadcasting messages to multiple subscribers across different server instances.
