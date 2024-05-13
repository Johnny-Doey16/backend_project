package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"

	"time"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/cache"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func sendPostToClient(posts []sqlc.GetPostsRecommendationRow, redisCache cache.Cache, stream pb.SocialMedia_PostStreamServer, uid uuid.UUID, postLikeKey string) error {
	for _, post := range posts {
		if err := stream.Send(&pb.PostStreamResponse{
			Data: &pb.PostStreamResponse_Post{
				Post: convertToPbPost(stream.Context(), redisCache, &post, uid, postLikeKey),
			},
		}); err != nil {
			return status.Errorf(codes.Internal, "error sending post: %s", err)
		}
	}
	return nil
}

func convertToPbPost(ctx context.Context, redisCache cache.Cache, post *sqlc.GetPostsRecommendationRow, uid uuid.UUID, postLikesKey string) *pb.Post {
	var imageUrls []string

	err := json.Unmarshal(post.PostImageUrls, &imageUrls)
	if err != nil {
		// Handle the error
		log.Printf("cannot unmarshal data: %s", err)
	}

	// Get like and likes from cache
	/*likeCountStr, _ := redisCache.GetKey(ctx, fmt.Sprintf("%s:%s", postLikesKey, post.ID))
	// Like in redis
	userLikesKey := fmt.Sprintf("%s:%s:users", postLikesKey, post.ID)
	// Check if the user has already liked the post
	liked, _ := redisCache.SIsMember(ctx, userLikesKey, uid.String())
	*/
	userLikesKey := fmt.Sprintf("%s:%s:users", postLikesKey, post.ID)
	liked, _ := redisCache.SIsMember(ctx, userLikesKey, uid.String())

	// Convert the sqlc.GetPostsRecommendationRow to the protobuf Post message
	var title *string

	return &pb.Post{
		UserId:       post.UserID.String(),
		Username:     post.Username.String,
		ProfileImage: post.UserImageUrl.String,
		Content:      post.Content,
		PostId:       post.ID.String(),
		ImageUrls:    imageUrls,
		Reason:       &post.Reason,
		Title:        InterfaceToStr(post.Title, title), // ! check what this is
		Metrics: &pb.PostMetrics{
			Likes:    post.Likes.Int32, // convertStrToInt32(likeCountStr), //
			Views:    post.Views.Int32,
			Comments: post.Comments.Int32,
			Reposts:  post.Reposts.Int32,
			Liked:    liked, // TODO: Get liked from db
		},
		Timestamp:        timestamppb.New(post.CreatedAt.Time),
		IsVerified:       post.IsVerified.Bool,
		FollowingsCounts: int64(post.FollowingCount.Int32),
		FollowersCounts:  int64(post.FollowersCount.Int32),
	}
}

func RetrievePosts(stream pb.SocialMedia_PostStreamServer, store *sqlc.Store, redisCache cache.Cache, uid uuid.UUID, postKey, postLikeKey string, postCacheExpr time.Duration, offset, limit int32) error {

	// Assuming a cache is in place, attempt to retrieve posts from cache first

	posts, err := store.GetPostsRecommendation(stream.Context(), sqlc.GetPostsRecommendationParams{
		UserID: uid,
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return status.Errorf(codes.Internal, "error fetching posts from db: %s", err)
	}

	// send post to client

	// ! Retrieve from cache before fetching from db
	/*
		posts, cacheErr := retrieveCachedPosts(redisCache, stream.Context(), postKey)
		if cacheErr != nil || len(posts) == 0 {

			// If cache retrieval fails or cache is empty, fallback to database
			var err error

			posts, err = store.GetPostsRecommendation(stream.Context(), sqlc.GetPostsRecommendationParams{
				UserID: uid,
				Offset: offset,
				Limit:  limit,
			})
			if err != nil {
				return status.Errorf(codes.Internal, "error fetching posts from db: %s", err)
			}
			// log.Println("Got From SQL", len(posts), " from the db")

			// After fetching from the database, add these posts to the cache
			cacheErr = CachePost(ctx, redisCache, posts, postKey, postCacheExpr)
			if cacheErr != nil {
				log.Printf("failed to add posts to cache: %s", cacheErr)
			}
		}
	*/

	if err := sendPostToClient(posts, redisCache, stream, uid, postLikeKey); err != nil {
		return status.Errorf(codes.Internal, "error sending post: %s", err)
	}
	return nil
}

func RefreshStream(stream pb.SocialMedia_PostStreamServer, store *sqlc.Store, redisCache cache.Cache, uid uuid.UUID, postChanKey, postLikeKey string) error {
	for {
		// Receive a PageRequest from the client
		postReq, err := stream.Recv()
		if err == io.EOF {
			// Client has closed the stream
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Internal, "Error receiving page request: %v", err)
		}

		// Pagination
		pageSize := postReq.PageSize
		pageNumber := postReq.PageNumber
		offset := (pageNumber - 1) * pageSize

		// Get refresh/next page for
		posts, err := store.GetPostsRecommendation(stream.Context(), sqlc.GetPostsRecommendationParams{
			UserID: uid,
			Offset: offset,
			Limit:  pageSize,
		})
		if err != nil {
			return status.Errorf(codes.Internal, "Error fetching posts from db: %v", err)
		}

		for _, post := range posts {
			if err := stream.Send(&pb.PostStreamResponse{
				Data: &pb.PostStreamResponse_Post{
					Post: convertToPbPost(stream.Context(), redisCache, &post, uid, postLikeKey),
				},
			}); err != nil {
				return status.Errorf(codes.Internal, "Error sending post to stream: %v", err)
			}
		}

		if stream.Context().Err() != nil {
			log.Println("RefreshStream: Context canceled. Exiting.")
			return nil
		}
	}
}

func SendLikesUpdatesToClient(redisCache cache.Cache, stream pb.SocialMedia_PostStreamServer, likeKey, postLikesKey string) {
	likesPubSub := redisCache.Subscribe(stream.Context(), likeKey)
	defer likesPubSub.Close()

	for {
		msg, err := likesPubSub.ReceiveMessage(stream.Context())
		if err != nil {
			if stream.Context().Err() != nil {
				return // Stream context was canceled (client disconnected)
			}
			log.Printf("Error receiving message: %v", err)
			continue
		}

		var likeUpdate LikeUpdate
		err = json.Unmarshal([]byte(msg.Payload), &likeUpdate)
		if err != nil {
			log.Printf("Error unmarshalling like update: %v", err)
			continue
		}

		// Here's where you construct and send the PostMetricsUpdate message
		var postMetricsUpdate pb.PostMetricsUpdate
		postMetricsUpdate.PostId = likeUpdate.PostID
		// Get the latest likes count from Redis
		likeCountKey := fmt.Sprintf("%s:%s", postLikesKey, likeUpdate.PostID)
		likeCount, err := redisCache.GetKey(stream.Context(), likeCountKey)
		if err != nil {
			// log.Printf("Error getting like count: %v", err)
			continue
		}

		// Convert to int32
		int32Value, err := strconv.ParseInt(likeCount, 10, 32)
		if err != nil {
			log.Printf("Error converting like count: %v", err)
			continue
		}

		postMetricsUpdate.Likes = int32(int32Value)

		// Send the metrics update to the client
		if err := stream.Send(&pb.PostStreamResponse{
			Data: &pb.PostStreamResponse_MetricsUpdate{
				MetricsUpdate: &postMetricsUpdate,
			},
		}); err != nil {
			log.Printf("Error sending metrics update to client: %v", err)
			return
		}

		// Update post like count in database or cache
		// err = s.updatePostLikes(stream.Context(), likeUpdate.PostID, likeUpdate.UserID, likeUpdate.Increment)
		// if err != nil {
		// 	// Handle the error appropriately
		// 	log.Printf("Error updating post likes: %v", err)
		// }
	}
}

func SendViewsUpdatesToClient(redisCache cache.Cache, stream pb.SocialMedia_PostStreamServer, viewKey string) {
	viewsPubSub := redisCache.Subscribe(stream.Context(), viewKey)
	defer viewsPubSub.Close()

	for {
		msg, err := viewsPubSub.ReceiveMessage(stream.Context())
		if err != nil {
			if stream.Context().Err() != nil {
				return // Stream context was canceled (client disconnected)
			}
			log.Printf("Error receiving message: %v", err)
			continue
		}

		var viewUpdate ViewUpdate
		err = json.Unmarshal([]byte(msg.Payload), &viewUpdate)
		if err != nil {
			log.Printf("Error unmarshalling like update: %v", err)
			continue
		}

		// Here's where you construct and send the PostMetricsUpdate message
		var postMetricsUpdate pb.PostMetricsUpdate
		postMetricsUpdate.PostId = viewUpdate.PostID
		// Get the latest likes count from Redis
		viewCountKey := fmt.Sprintf("post_views:%s", viewUpdate.PostID)
		viewCount, err := redisCache.GetKey(stream.Context(), viewCountKey)
		if err != nil {
			// log.Printf("Error getting like count: %v", err)
			continue
		}

		// Convert to int32
		int32Value, err := strconv.ParseInt(viewCount, 10, 32)
		if err != nil {
			log.Printf("Error converting like count: %v", err)
			continue
		}

		postMetricsUpdate.Views = int32(int32Value)

		// Send the metrics update to the client
		if err := stream.Send(&pb.PostStreamResponse{
			Data: &pb.PostStreamResponse_MetricsUpdate{
				MetricsUpdate: &postMetricsUpdate,
			},
		}); err != nil {
			log.Printf("Error sending metrics update to client: %v", err)
			return
		}
	}
}

func SendCommentsUpdatesToClient(redisCache cache.Cache, stream pb.SocialMedia_PostStreamServer, commentKey, postCommentsKey string) {
	commentPubSub := redisCache.Subscribe(stream.Context(), commentKey)
	defer commentPubSub.Close()

	for {
		msg, err := commentPubSub.ReceiveMessage(stream.Context())
		if err != nil {
			if stream.Context().Err() != nil {
				return // Stream context was canceled (client disconnected)
			}
			log.Printf("Error receiving message: %v", err)
			continue
		}

		var commentUpdate CommentUpdate
		err = json.Unmarshal([]byte(msg.Payload), &commentUpdate)
		if err != nil {
			log.Printf("Error unmarshalling comment update: %v", err)
			continue
		}

		// Here's where you construct and send the PostMetricsUpdate message
		var postMetricsUpdate pb.PostMetricsUpdate
		postMetricsUpdate.PostId = commentUpdate.PostID

		// Get the latest likes count from Redis
		commentCountKey := fmt.Sprintf("%s:%s", postCommentsKey, commentUpdate.PostID)
		commentCount, err := redisCache.GetKey(stream.Context(), commentCountKey)
		if err != nil {
			// log.Printf("Error getting like count: %v", err)
			continue
		}

		// Convert to int32
		int32Value, err := strconv.ParseInt(commentCount, 10, 32)
		if err != nil {
			log.Printf("Error converting like count: %v", err)
			continue
		}

		postMetricsUpdate.Comments = int32(int32Value)

		// Send the metrics update to the client
		if err := stream.Send(&pb.PostStreamResponse{
			Data: &pb.PostStreamResponse_MetricsUpdate{
				MetricsUpdate: &postMetricsUpdate,
			},
		}); err != nil {
			log.Printf("Error sending metrics update to client: %v", err)
			return
		}
	}
}

func SendRepostsUpdatesToClient(redisCache cache.Cache, stream pb.SocialMedia_PostStreamServer, repostChannelKey, repostsKey string) {
	repostPubSub := redisCache.Subscribe(stream.Context(), repostChannelKey)
	defer repostPubSub.Close()

	for {
		msg, err := repostPubSub.ReceiveMessage(stream.Context())
		if err != nil {
			if stream.Context().Err() != nil {
				return // Stream context was canceled (client disconnected)
			}
			log.Printf("Error receiving message: %v", err)
			continue
		}

		var repostUpdate RepostUpdate
		err = json.Unmarshal([]byte(msg.Payload), &repostUpdate)
		if err != nil {
			log.Printf("Error unmarshalling comment update: %v", err)
			continue
		}

		// Here's where you construct and send the PostMetricsUpdate message
		var postMetricsUpdate pb.PostMetricsUpdate
		postMetricsUpdate.PostId = repostUpdate.PostID

		// Get the latest likes count from Redis
		repostCountKey := fmt.Sprintf("%s:%s", repostsKey, repostUpdate.PostID)
		repostCount, err := redisCache.GetKey(stream.Context(), repostCountKey)
		if err != nil {
			// log.Printf("Error getting like count: %v", err)
			continue
		}

		// Convert to int32
		int32Value, err := strconv.ParseInt(repostCount, 10, 32)
		if err != nil {
			log.Printf("Error converting like count: %v", err)
			continue
		}

		postMetricsUpdate.Reposts = int32(int32Value)

		// Send the metrics update to the client
		if err := stream.Send(&pb.PostStreamResponse{
			Data: &pb.PostStreamResponse_MetricsUpdate{
				MetricsUpdate: &postMetricsUpdate,
			},
		}); err != nil {
			log.Printf("Error sending metrics update to client: %v", err)
			return
		}
	}
}

func retrieveCachedPosts(cacheClient cache.Cache, ctx context.Context, postKey string) ([]sqlc.GetPostsRecommendationRow, error) {
	//fmt.Sprintf("posts:%s", postKey)
	// Assuming cacheClient is a Redis client
	cachedPosts, cacheErr := cacheClient.GetKey(ctx, postKey) //.Result()
	if cacheErr != nil {
		return nil, cacheErr
	}

	var posts []sqlc.GetPostsRecommendationRow
	err := json.Unmarshal([]byte(cachedPosts), &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
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

func processPageRequest(stream pb.SocialMedia_PostStreamServer, store *sqlc.Store, redisCache cache.Cache, postReq *pb.PostRequest, uid uuid.UUID, postLikeKey string) {
	// Pagination
	pageSize := postReq.PageSize
	pageNumber := postReq.PageNumber
	offset := (pageNumber - 1) * pageSize

	// Get refresh/next page for posts
	posts, err := store.GetPostsRecommendation(stream.Context(), sqlc.GetPostsRecommendationParams{
		UserID: uid,
		Offset: offset,
		Limit:  pageSize,
	})
	if err != nil {
		log.Printf("error fetching posts from db: %s", err)
		return
	}

	for _, post := range posts {
		pbPost := convertToPbPost(stream.Context(), redisCache, &post, uid, postLikeKey)
		if err := stream.Send(&pb.PostStreamResponse{
			Data: &pb.PostStreamResponse_Post{
				Post: pbPost,
			},
		}); err != nil {
			log.Printf("1 Error sending post to stream: %v", err)
			return
		}
	}
}

// * Depreciated functions
func IsPostRecommended(store *sqlc.Store, ctx context.Context, uid uuid.UUID) bool {
	shouldRecommend, err := store.IsPostRecommended(ctx, uid)
	if err != nil {
		return false
	}
	log.Println("Should recommend", shouldRecommend)
	return shouldRecommend
}

/**
func isFollowing(store *sqlc.Store, ctx context.Context, uid uuid.UUID, authorID string) bool {
	author, _ := services.StrToUUID(authorID)
	_, err := store.IsUserFollowing(ctx, sqlc.IsUserFollowingParams{
		FollowerUserID:  uid,
		FollowingUserID: author,
	})
	if err != nil {
		return false
	}
	return true
}
*/

// For send new posts to client
func processNewPost(stream pb.SocialMedia_PostStreamServer, store *sqlc.Store, post *pb.Post, uid uuid.UUID, postLikeKey string) {
	// Check if the current user is following the author of the post
	if IsPostRecommended(store, stream.Context(), uid) {
		if err := stream.Send(&pb.PostStreamResponse{
			Data: &pb.PostStreamResponse_Post{
				Post: post,
			},
		}); err != nil {
			log.Printf("Error sending post to stream: %v", err)
		}
	}
}

func ListenForRefresh(stream pb.SocialMedia_PostStreamServer, store *sqlc.Store, redisCache cache.Cache, uid uuid.UUID, postChanKey, postLikeKey string) {
	for {
		// Receive a PageRequest from the client
		postReq, err := stream.Recv()
		if err == io.EOF {
			// Client has closed the stream
			return
		}
		if err != nil {
			log.Printf("Error receiving page request: %v", err)
			return
		}
		log.Println("Calling REQ with", postReq)

		// Pagination
		pageSize := postReq.PageSize
		pageNumber := postReq.PageNumber
		offset := (pageNumber - 1) * pageSize

		// Get refresh/next page for
		posts, err := store.GetPostsRecommendation(stream.Context(), sqlc.GetPostsRecommendationParams{
			UserID: uid,
			Offset: offset,
			Limit:  pageSize,
		})
		if err != nil {
			log.Printf("error fetching posts from db: %s", err)
			return
		}

		for _, post := range posts {
			if err := stream.Send(&pb.PostStreamResponse{
				Data: &pb.PostStreamResponse_Post{
					Post: convertToPbPost(stream.Context(), redisCache, &post, uid, postLikeKey),
				},
			}); err != nil {
				log.Printf("Error sending post to stream: %v", err)
				return
			}
		}

		if stream.Context().Err() != nil {
			log.Println("ListenForRefresh: Context canceled. Exiting.")
			return
		}

	}
}

func ListenForRefresh2(stream pb.SocialMedia_PostStreamServer, store *sqlc.Store, redisCache cache.Cache, uid uuid.UUID, postChanKey, postLikeKey string) {
	log.Println("ListenForRefresh: Started")
	defer log.Println("ListenForRefresh: Stopped")

	for {
		log.Println("ListenForRefresh: Waiting for page request")

		// Receive a PageRequest from the client
		postReq, err := stream.Recv()
		if err == io.EOF {
			// Client has closed the stream
			log.Println("ListenForRefresh: Client closed the stream")
			return
		}
		if err != nil {
			log.Printf("ListenForRefresh: Error receiving page request: %v", err)
			return
		}

		log.Printf("ListenForRefresh: Received page request: %+v", postReq)

		// Pagination
		pageSize := postReq.PageSize
		pageNumber := postReq.PageNumber
		offset := (pageNumber - 1) * pageSize

		log.Printf("ListenForRefresh: Fetching posts from offset %d, limit %d", offset, pageSize)

		// Get refresh/next page for
		posts, err := store.GetPostsRecommendation(stream.Context(), sqlc.GetPostsRecommendationParams{
			UserID: uid,
			Offset: offset,
			Limit:  pageSize,
		})
		if err != nil {
			log.Printf("ListenForRefresh: Error fetching posts from db: %s", err)
			return
		}

		log.Printf("ListenForRefresh: Fetched %d posts", len(posts))

		for _, post := range posts {
			log.Printf("ListenForRefresh: Sending post to client: %s", post.ID)

			// Lock the mutex before sending updates to the client
			if err := stream.Send(&pb.PostStreamResponse{
				Data: &pb.PostStreamResponse_Post{
					Post: convertToPbPost(stream.Context(), redisCache, &post, uid, postLikeKey),
				},
			}); err != nil {
				log.Printf("ListenForRefresh: Error sending post to stream: %v", err)
				return
			}
		}

		// Ensure the context is still active before waiting for the next page request
		if stream.Context().Err() != nil {
			log.Println("ListenForRefresh: Context canceled. Exiting.")
			return
		}

		log.Println("ListenForRefresh: Waiting for the next page request")
	}
}

func SendNewPostsToClientOld(stream pb.SocialMedia_PostStreamServer, store *sqlc.Store, redisCache cache.Cache, uid uuid.UUID, postChanKey, postLikeKey string) {
	// Subscribe to the Redis channel
	pubsub := redisCache.Subscribe(stream.Context(), postChanKey) //PostChannelKey
	defer pubsub.Close()

	for {

		// ! Listen for new posts
		msg, err := pubsub.ReceiveMessage(stream.Context())
		if err != nil {
			if stream.Context().Err() != nil {
				return // Stream context was canceled (client disconnected)
			}
			log.Printf("Error receiving message: %v", err)
			continue
		}
		log.Println("Receieved", msg.Payload)

		var post pb.Post
		err = json.Unmarshal([]byte(msg.Payload), &post)
		if err != nil {
			log.Printf("Error unmarshalling post: %v", err)
			continue
		}

		// Check if the current user is following the author of the post
		// Note: Assuming `isPostRecommended` checks if the post should be sent to the user.
		if IsPostRecommended(store, stream.Context(), uid) {

			if err := stream.Send(&pb.PostStreamResponse{
				Data: &pb.PostStreamResponse_Post{
					Post: &post,
				},
			}); err != nil {
				log.Printf("Error sending post to stream: %v", err)
				return
			}
		}
	}
}

func InterfaceToStr(strInterface interface{}, content *string) *string {
	nilStr := ""
	if strInterface != nil {
		titleStr, ok := strInterface.(string)
		if ok {
			content = &titleStr
		}
	} else {
		content = &nilStr
	}
	return content
}

/*func SendNewPostsToClient1(stream pb.SocialMedia_PostStreamServer, store *sqlc.Store, redisCache cache.Cache, uid uuid.UUID, postChanKey, postLikeKey string) {
	// Subscribe to the Redis channel
	log.Printf("Subscribing to Redis channel: %s", postChanKey)
	pubsub := redisCache.Subscribe(stream.Context(), postChanKey)
	defer pubsub.Close()

	// Channel to handle posts from Redis
	newPostChan := make(chan *pb.Post)

	// Goroutine to listen to new posts from Redis
	go func() {
		for {
			msg, err := pubsub.ReceiveMessage(stream.Context())
			if err != nil {
				if stream.Context().Err() != nil {
					log.Printf("Stream context was canceled, stopping Redis listener goroutine")
					return // Stream context was canceled (client disconnected)
				}
				log.Printf("Error receiving message from Redis: %v", err)
				continue
			}

			var post pb.Post
			err = json.Unmarshal([]byte(msg.Payload), &post)
			if err != nil {
				log.Printf("Error unmarshalling post from Redis message: %v", err)
				continue
			}

			log.Printf("New post received from Redis: %v", post)
			newPostChan <- &post
		}
	}()

	log.Printf("Listening for new posts and page requests...")
	for {
		select {
		case <-stream.Context().Done():
			log.Printf("Stream context was canceled, exiting SendNewPostsToClient function")
			return // Stream context was canceled (client disconnected)
		case postReq, ok := <-newPostChan:
			if !ok {
				log.Printf("New posts channel was closed, exiting SendNewPostsToClient function")
				return // Channel was closed, exit the loop
			}
			log.Printf("Processing new post...")
			processNewPost(stream, store, postReq, uid, postLikeKey)
			// Remove the default case to prevent blocking on stream.Recv() when waiting for new posts
		}
	}
}

func SendNewPostsToClient2(stream pb.SocialMedia_PostStreamServer, store *sqlc.Store, redisCache cache.Cache, uid uuid.UUID, postChanKey, postLikeKey string) {
	// Subscribe to the Redis channel
	log.Printf("Subscribing to Redis channel: %s", postChanKey)
	pubsub := redisCache.Subscribe(stream.Context(), postChanKey)
	defer pubsub.Close()

	// Channel to handle posts from Redis
	newPostChan := make(chan *pb.Post)
	newPostListChan := make(chan *[]sqlc.GetPostsRecommendationRow)

	// Goroutine to listen to new posts from Redis
	go func() {
		for {
			msg, err := pubsub.ReceiveMessage(stream.Context())
			if err != nil {
				if stream.Context().Err() != nil {
					log.Printf("Stream context was canceled, stopping Redis listener goroutine")
					return // Stream context was canceled (client disconnected)
				}
				log.Printf("Error receiving message from Redis: %v", err)
				continue
			}

			var post pb.Post
			err = json.Unmarshal([]byte(msg.Payload), &post)
			if err != nil {
				log.Printf("Error unmarshalling post from Redis message: %v", err)
				continue
			}

			log.Printf("New post received from Redis: %v", post)
			newPostChan <- &post
		}
	}()

	// Go routine to paginate posts
	go func() {
		postReq, err := stream.Recv()
		if err == io.EOF {
			// Client has closed the stream
			return
		}
		if err != nil {
			log.Printf("Error receiving page request: %v", err)
			// continue
		}
		pageSize := postReq.PageSize
		pageNumber := postReq.PageNumber
		offset := (pageNumber - 1) * pageSize

		// Get refresh/next page for posts
		posts, err := store.GetPostsRecommendation(stream.Context(), sqlc.GetPostsRecommendationParams{
			UserID: uid,
			Offset: offset,
			Limit:  pageSize,
		})
		if err != nil {
			log.Printf("error fetching posts from db: %s", err)
			return
		}
		newPostListChan <- &posts

	}()

	log.Printf("Listening for new posts and page requests...")
	for {
		select {
		case <-stream.Context().Done():
			log.Printf("Stream context was canceled, exiting SendNewPostsToClient function")
			return // Stream context was canceled (client disconnected)
		case postReq, ok := <-newPostChan:
			if !ok {
				log.Printf("New posts channel was closed, exiting SendNewPostsToClient function")
				return // Channel was closed, exit the loop
			}
			log.Printf("Processing new post...")
			processNewPost(stream, store, postReq, uid, postLikeKey)
			// Remove the default case to prevent blocking on stream.Recv() when waiting for new posts

			// case for when there is data in stream.Recv()
		case posts, ok := <-newPostListChan:
			if !ok {
				log.Printf("New posts channel was closed, exiting SendNewPostsToClient function")
				return // Channel was closed, exit the loop
			}
			for _, post := range *posts {
				pbPost := convertToPbPost(stream.Context(), redisCache, &post, uid, postLikeKey)
				if err := stream.Send(&pb.PostStreamResponse{
					Data: &pb.PostStreamResponse_Post{
						Post: pbPost,
					},
				}); err != nil {
					log.Printf("1 Error sending post to stream: %v", err)
					return
				}
			}
		}
	}
}
*/

func SendNewPostsToClient(stream pb.SocialMedia_PostStreamServer, store *sqlc.Store, redisCache cache.Cache, uid uuid.UUID, postChanKey, postLikeKey string) {
	// Subscribe to the Redis channel
	pubsub := redisCache.Subscribe(stream.Context(), postChanKey)
	defer pubsub.Close()

	// Channel to handle posts from Redis
	newPostChan := make(chan *pb.Post)

	// Goroutine to listen to new posts from Redis
	go func() {
		for {
			msg, err := pubsub.ReceiveMessage(stream.Context())
			if err != nil {
				if stream.Context().Err() != nil {
					return // Stream context was canceled (client disconnected)
				}
				log.Printf("Error receiving message: %v", err)
				continue
			}

			var post pb.Post
			err = json.Unmarshal([]byte(msg.Payload), &post)
			if err != nil {
				log.Printf("Error unmarshalling post: %v", err)
				continue
			}

			newPostChan <- &post
		}
	}()

	for {
		select {
		case <-stream.Context().Done():
			return // Stream context was canceled (client disconnected)
		case postReq, ok := <-newPostChan:
			if !ok {
				return // Channel was closed, exit the loop
			}
			// Use a function to process new posts
			processNewPost(stream, store, postReq, uid, postLikeKey)
		default:
			// Receive a PageRequest from the client
			postReq, err := stream.Recv()
			if err == io.EOF {
				// Client has closed the stream
				return
			}
			if err != nil {
				log.Printf("Error receiving page request: %v", err)
				continue
			}
			// Use a function to process page requests
			processPageRequest(stream, store, redisCache, postReq, uid, postLikeKey)
		}
	}
}
