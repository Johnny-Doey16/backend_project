package gapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/steve-mir/diivix_backend/cache"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"github.com/steve-mir/diivix_backend/utils"
	"github.com/steve-mir/diivix_backend/worker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO: 3. Consider replacing encoding/json with jsoniter or any faster json encoder
// Testing: Conduct load testing and optimize based on the results to ensure the system can handle the expected traffic.
// Caching Layer: The caching strategy should be optimized for scalability. Consider using a distributed cache and setting appropriate cache expiration times.
// Concurrency and Throttling: Ensure that the server can handle multiple concurrent streams without degrading performance. Implement throttling if necessary to prevent overloading the server.
// Database Access: Optimize database queries to be efficient and consider using read replicas for fetching data to distribute the load.
// Code Optimization: Review and refactor the code for computational efficiency. For example, the isPostRecommended function might be a performance bottleneck if it's not optimized for large-scale data.
// Resilience: Implement retries with exponential backoff and circuit breakers to make the system more resilient to failures.

const (
	defaultPageNumber = 1
	defaultPageSize   = 10
)

func (s *SocialMediaServer) PostStream(stream pb.SocialMedia_PostStreamServer) error {
	sTime := utils.StartMemCal()
	claims, ok := stream.Context().Value("payloadKey").(*token.Payload)
	if !ok {
		return status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}
	// userId, err := services.StrToUUID("558c60aa-977f-4d38-885b-e813656371ac")
	// if err != nil {
	// 	return status.Errorf(codes.Internal, "parsing uid: %s", err)
	// }

	// Calculate the offset based on the page number and page size.
	offset := (defaultPageNumber - 1) * defaultPageSize

	// Retrieve posts either from cache or database
	if err := services.RetrievePosts(stream, s.store, s.redisCache, claims.UserId, PostCacheKey, PostLikesKey, PostCacheExpr, int32(offset), int32(defaultPageSize)); err != nil {
		return err
	}

	/* // !
	// Subscribe to the Redis channel for real-time updates
	pubsub := s.redisCache.Subscribe(stream.Context(), PostChannelKey)
	defer pubsub.Close()

	// Handle errors from subscription
	_, err = pubsub.Receive(stream.Context())
	if err != nil {
		return status.Errorf(codes.Internal, "error subscribing to Redis: %s", err)
	}
	*/

	// Start a goroutine to listen for new posts from Redis Pub/Sub
	go func() {
		err := services.RefreshStream(stream, s.store, s.redisCache, claims.UserId, PostChannelKey, PostLikesKey)
		if err != nil {
			// Handle the error appropriately, e.g., log it or perform cleanup
			log.Printf("Error in RefreshStream: %v", err)
		}
	}()

	// Initialize and run the worker pool
	workerPool := worker.NewWorkerPool(10) // Set the number of workers based on your needs
	workerPool.Run()
	defer workerPool.Shutdown()
	// Create a channel for update messages
	updateChan := make(chan *pb.PostStreamResponse)

	go s.subscribeToRedisChannels(stream.Context(), updateChan, s.redisCache, workerPool, PostChannelKey, LikeChannelKey, CommentChannelKey, RepostChannelKey, ViewChannelKey)
	// Start a single goroutine to send updates from updateChan to the client
	go func() {
		for update := range updateChan {
			if err := stream.Send(update); err != nil {
				log.Printf("Error sending update to stream: %v", err)
				return
			}
		}
	}()

	// Wait for the client to disconnect
	<-stream.Context().Done()
	utils.EndMemCal(sTime)
	return stream.Context().Err()
}

func (s *SocialMediaServer) subscribeToRedisChannels(ctx context.Context, updateChan chan<- *pb.PostStreamResponse, redisCache cache.Cache, workerPool *worker.WorkerPool, channels ...string) {
	for _, channel := range channels {
		pubsub := redisCache.Subscribe(ctx, channel)
		// Pass the channel as an argument to the closure to ensure it captures the current value.
		workerPool.AddJob(func(ch string) func() {
			return func() {
				defer pubsub.Close() // Defer the closing of the subscription outside the loop
				for {
					msg, err := pubsub.ReceiveMessage(ctx)
					if err != nil {
						if ctx.Err() != nil {
							return // Context was canceled (client disconnected)
						}
						log.Printf("Error receiving message from channel %s: %v", ch, err)
						return
					}

					// Process the message and send to updateChan
					log.Println("Received message on channel", ch)
					update, err := s.processReceivedMessage(ctx, redisCache, msg, ch)
					if err != nil {
						log.Printf("Error processing message from channel %s: %v", ch, err)
						continue
					}
					updateChan <- update
				}
			}
		}(channel)) // Pass channel as an argument to the closure
	}
}

func (s *SocialMediaServer) processReceivedMessage(ctx context.Context, redisCache cache.Cache, msg *redis.Message, channel string) (*pb.PostStreamResponse, error) {
	log.Printf("Inside Process Rec Msg %s", channel)

	switch channel {
	case PostChannelKey:
		return processPostMessage(msg)
	case LikeChannelKey:
		return processLikeMessage(ctx, redisCache, msg)
	case CommentChannelKey:
		return processCommentMessage(ctx, redisCache, msg)
	case RepostChannelKey:
		return processRepostMessage(ctx, redisCache, msg)
	case ViewChannelKey:
		return nil, nil // processViewMessage(msg.Payload)
	default:
		return nil, fmt.Errorf("unknown channel: %s", channel)
	}
}

func processPostMessage(msg *redis.Message) (*pb.PostStreamResponse, error) {
	var post pb.Post
	if err := json.Unmarshal([]byte(msg.Payload), &post); err != nil {
		return nil, fmt.Errorf("error unmarshalling new post: %v", err)
	}
	return &pb.PostStreamResponse{
		Data: &pb.PostStreamResponse_Post{Post: &post},
	}, nil
}

func processLikeMessage(ctx context.Context, redisCache cache.Cache, msg *redis.Message) (*pb.PostStreamResponse, error) {
	var likeUpdate services.LikeUpdate
	if err := json.Unmarshal([]byte(msg.Payload), &likeUpdate); err != nil {
		log.Printf("Error unmarshalling like update: %v", err)
		return nil, nil
	}

	postMetricsUpdate, err := buildPostMetricsUpdateForLikes(ctx, redisCache, likeUpdate.PostID)
	if err != nil {
		return nil, err
	}

	return &pb.PostStreamResponse{
		Data: &pb.PostStreamResponse_MetricsUpdate{
			MetricsUpdate: postMetricsUpdate,
		},
	}, nil
}

func processCommentMessage(ctx context.Context, redisCache cache.Cache, msg *redis.Message) (*pb.PostStreamResponse, error) {
	var commentUpdate services.CommentUpdate
	if err := json.Unmarshal([]byte(msg.Payload), &commentUpdate); err != nil {
		log.Printf("Error unmarshalling comment update: %v", err)
		return nil, nil
	}

	postMetricsUpdate, err := buildPostMetricsUpdateForComments(ctx, redisCache, commentUpdate.PostID)
	if err != nil {
		return nil, err
	}

	return &pb.PostStreamResponse{
		Data: &pb.PostStreamResponse_MetricsUpdate{
			MetricsUpdate: postMetricsUpdate,
		},
	}, nil
}

func processRepostMessage(ctx context.Context, redisCache cache.Cache, msg *redis.Message) (*pb.PostStreamResponse, error) {
	var repostUpdate services.RepostUpdate
	if err := json.Unmarshal([]byte(msg.Payload), &repostUpdate); err != nil {
		log.Printf("Error unmarshalling repost update: %v", err)
		return nil, nil
	}

	postMetricsUpdate, err := buildPostMetricsUpdateForRepost(ctx, redisCache, repostUpdate.PostID)
	if err != nil {
		return nil, err
	}

	return &pb.PostStreamResponse{
		Data: &pb.PostStreamResponse_MetricsUpdate{
			MetricsUpdate: postMetricsUpdate,
		},
	}, nil
}

func buildPostMetricsUpdateForLikes(ctx context.Context, redisCache cache.Cache, postID string) (*pb.PostMetricsUpdate, error) {
	var postMetricsUpdate pb.PostMetricsUpdate
	postMetricsUpdate.PostId = postID

	likeCountKey := fmt.Sprintf("%s:%s", PostLikesKey, postID)
	likeCount, err := redisCache.GetKey(ctx, likeCountKey)
	if err != nil {
		log.Printf("Error getting like count: %v", err)
		return nil, err
	}

	int32Value, err := strconv.ParseInt(likeCount, 10, 32)
	if err != nil {
		log.Printf("Error converting like count: %v", err)
		return nil, err
	}

	postMetricsUpdate.Likes = int32(int32Value)
	return &postMetricsUpdate, nil
}

func buildPostMetricsUpdateForComments(ctx context.Context, redisCache cache.Cache, postID string) (*pb.PostMetricsUpdate, error) {
	var postMetricsUpdate pb.PostMetricsUpdate
	postMetricsUpdate.PostId = postID

	commentCountKey := fmt.Sprintf("%s:%s", PostCommentsKey, postID)
	commentCount, err := redisCache.GetKey(ctx, commentCountKey)
	if err != nil {
		log.Printf("Error getting comments count: %v", err)
		return nil, err
	}

	int32Value, err := strconv.ParseInt(commentCount, 10, 32)
	if err != nil {
		log.Printf("Error converting comment count: %v", err)
		return nil, err
	}

	postMetricsUpdate.Comments = int32(int32Value)
	return &postMetricsUpdate, nil
}

func buildPostMetricsUpdateForRepost(ctx context.Context, redisCache cache.Cache, postID string) (*pb.PostMetricsUpdate, error) {
	var postMetricsUpdate pb.PostMetricsUpdate
	postMetricsUpdate.PostId = postID

	repostsCountKey := fmt.Sprintf("%s:%s", RepostsKey, postID)
	repostsCount, err := redisCache.GetKey(ctx, repostsCountKey)
	if err != nil {
		log.Printf("Error getting reposts count: %v", err)
		return nil, err
	}

	int32Value, err := strconv.ParseInt(repostsCount, 10, 32)
	if err != nil {
		log.Printf("Error converting reposts count: %v", err)
		return nil, err
	}

	postMetricsUpdate.Reposts = int32(int32Value)
	return &postMetricsUpdate, nil
}

// Similar functions for other message types (e.g, processViewMessage)

/*
	// Measure time taken by encoding/json
	startTimeJSON := time.Now()
	var likeUpdate services.LikeUpdate
	err := json.Unmarshal([]byte(msg.Payload), &likeUpdate)
	if err != nil {
		log.Printf("Error unmarshalling post with encoding/json: %v", err)
	}
	log.Printf("Time taken by encoding/json: %v", time.Since(startTimeJSON))

	// Measure time taken by json-iterator
	startTimeJSONIterator := time.Now()
	var postJSONIterator services.LikeUpdate
	err = jsoniter.Unmarshal([]byte(msg.Payload), &postJSONIterator)
	if err != nil {
		log.Printf("Error unmarshalling post with json-iterator: %v", err)
	}
	log.Printf("Time taken by json-iterator: %v", time.Since(startTimeJSONIterator))
*/
