package gapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/cache"
	"github.com/steve-mir/diivix_backend/constants"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *SocialMediaServer) SendPostComment(ctx context.Context, req *pb.CommentRequest) (*pb.CommentResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)

	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}
	// commenterID, err := services.StrToUUID("e7679d8b-0eac-4ea2-93cd-0018ab995922")
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "parsing uid: %s", err)
	// }

	postID, postUserID, err := getIdsFromString(req.GetPostId(), req.GetPostUserId()) //services.StrToUUID(req.GetPostId())
	if err != nil {
		return nil, err
	}

	err = services.CheckComment(req.GetText())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "comment error %v", err)
	}

	// Create transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transaction error: %v", err.Error())
	}
	defer tx.Rollback()
	qtx := s.store.WithTx(tx)

	comment_id, t_comments, err := services.RunPostCommentConcurrent(tx, qtx, ctx, req.GetText(), postID, claims.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Send notification
	if postUserID != claims.UserId {
		services.SendNotification(s.taskDistributor, ctx, claims.UserId, postID, []uuid.UUID{postUserID}, constants.NotificationPostComment, "Comment", "", fmt.Sprintf("%s Commented on your post", claims.Username), "")
	}

	// Add comment count to redis
	err = addCommentCountToRedis(s.redisCache, ctx, PostCommentsKey, req.GetPostId(), t_comments)
	if err != nil {
		return nil, err
	}

	// broadcast comment event here
	if err := broadcastCommentEvent(ctx, s.redisCache, req.GetPostId(), claims.UserId); err != nil {
		return nil, err
	}

	return &pb.CommentResponse{
		CommentId: comment_id,
	}, nil
}

func getIdsFromString(postId, postUserId string) (uuid.UUID, uuid.UUID, error) {
	postID, err := services.StrToUUID(postId)
	if err != nil {
		return uuid.Nil, uuid.Nil, status.Errorf(codes.Internal, "error generating post id %v", err)
	}

	postUserID, err := services.StrToUUID(postUserId)
	if err != nil {
		return uuid.Nil, uuid.Nil, status.Errorf(codes.Internal, "error generating post user id %v", err)
	}
	return postID, postUserID, nil
}

func addCommentCountToRedis(redisCache cache.Cache, ctx context.Context, key, postID string, counts int32) error {
	// Set the like count in Redis to match the database count
	commentCountKey := fmt.Sprintf("%s:%s", key, postID)
	err := redisCache.SetKey(ctx, commentCountKey, counts, time.Hour*24)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to set like count in Redis: %v", err.Error())
	}
	return nil
}

func broadcastCommentEvent(ctx context.Context, redisCache cache.Cache,
	postId string, uid uuid.UUID) error {
	commentUpdate := services.CommentUpdate{
		PostID: postId,
		UserID: uid.String(),
		// Additional fields as needed
	}
	commentUpdateBytes, err := json.Marshal(commentUpdate)
	if err != nil {
		log.Printf("Error marshalling comment update: %v", err)
		return err
	}

	// Publish the like update to the Redis channel
	err = redisCache.Publish(ctx, CommentChannelKey, commentUpdateBytes).Err()
	if err != nil {
		log.Printf("Error publishing like update: %v", err)
		return err
	}
	return nil
}
