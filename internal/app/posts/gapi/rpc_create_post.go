package gapi

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	PostCacheExpr = 24 * time.Hour
)

func (s *SocialMediaServer) broadcast(ctx context.Context, post *pb.Post) {
	// Serialize the post object to JSON for publishing.
	postBytes, err := json.Marshal(post)
	if err != nil {
		log.Printf("Error marshalling post: %v", err)
		return
	}

	// Publish the serialized post to the Redis channel.
	err = s.redisCache.Publish(ctx, PostChannelKey, postBytes).Err()
	if err != nil {
		log.Printf("Error publishing post: %v", err)
	}
}

// TODO: Check if image and texts (comment and post texts) are safe to upload
// Concurrency and Locking: The use of s.mu.Lock() suggests some shared mutable state. Investigate if this lock is truly necessary and consider using finer-grained locking or lock-free structures to reduce contention.
// Resource Utilization: Serializing and broadcasting a post is done within the request handling method CreatePost. This can be offloaded to a separate worker or service to return a response to the user faster and process the broadcast asynchronously.
// Caching Strategies: When caching the post, use a unique key based on the post ID to prevent collisions. Also, consider setting varying expiration times based on the popularity of the post.
// Message Publishing: When publishing to the Redis channel, ensure that you're not blocking the main thread and handle backpressure appropriately.
func (s *SocialMediaServer) CreatePost(stream pb.SocialMedia_CreatePostServer) error {
	ctx, cancel := context.WithTimeout(stream.Context(), time.Second*10) // Adjust the timeout as needed
	defer cancel()

	// TODO: Check whether to remove or leave
	s.mu.Lock()
	defer s.mu.Unlock()
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}
	// userId, err := services.StrToUUID("558c60aa-977f-4d38-885b-e813656371ac")
	// if err != nil {
	// 	return status.Errorf(codes.Internal, "parsing uid: %s", err)
	// }

	imageUrls, captions, content, err := services.GetCreatePostStreamData(stream, s.imageStore)
	if err != nil {
		return err
	}

	postID, sqlPost, err := services.CreatePostInDB(ctx, &s.taskDistributor, s.db, s.store, claims.UserId, content, claims.Username, imageUrls, captions)
	if err != nil {
		return err
	}

	s.broadcastNewPost(ctx, sqlPost, imageUrls)

	// Send the response back to the client
	return stream.SendAndClose(&pb.CreatePostResponse{
		PostId:    postID.String(),
		ImageUrls: imageUrls,
	})
	// return nil
}

func (s *SocialMediaServer) broadcastNewPost(ctx context.Context, sqlPost sqlc.Post, imageUrls []string) {
	// Create the gRPC response
	post := pb.Post{
		PostId:    sqlPost.ID.String(),
		Content:   sqlPost.Content,
		UserId:    sqlPost.UserID.String(),
		ImageUrls: imageUrls,
		Timestamp: timestamppb.New(sqlPost.CreatedAt.Time),
	}

	if err := services.CachePost(ctx, s.redisCache, &post, PostCacheKey, PostCacheExpr); err != nil {
		// Handle the error but you might not want to fail the whole operation
		// because the main action (inserting into DB) was successful.
		// Log this error instead of returning it.
		log.Printf("Failed to cache the post: %v\n", err)
	}

	// Broadcast the post to subscribers
	s.broadcast(ctx, &post)
}
