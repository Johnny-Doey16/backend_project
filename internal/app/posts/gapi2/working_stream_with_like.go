package gapi2

// func (s *SocialMediaServer) PostStream(stream pb.SocialMedia_PostStreamServer) error {
// 	// _, ok := stream.Context().Value("payloadKey").(*token.Payload)
// 	//
// 	//	if !ok {
// 	//		return status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
// 	//	}
// 	uid, _ := services.StrToUUID("ad7a61cd-14c5-4bbe-a3fe-0abdb585898a")

// 	// Retrieve posts either from cache or database
// 	posts, err := retrievePosts(stream, s.redisCache, s.store, uid)
// 	if err != nil {
// 		return err
// 	}

// 	// Send the retrieved posts to the client first
// 	for _, post := range posts {
// 		if err := stream.Send(convertToPbPost(&post)); err != nil {
// 			return status.Errorf(codes.Internal, "error sending post: %s", err)
// 		}
// 	}

// 	// Subscribe to the Redis channel for real-time updates
// 	pubsub := s.redisCache.Subscribe(stream.Context(), PostChannelKey)
// 	defer pubsub.Close()

// 	// Handle errors from subscription
// 	_, err = pubsub.Receive(stream.Context())
// 	if err != nil {
// 		return status.Errorf(codes.Internal, "error subscribing to Redis: %s", err)
// 	}

// 	// Start a goroutine to listen for new posts from Redis Pub/Sub
// 	go s.sendNewPostsToClient(stream, uid)

// 	// go s.sendLikesUpdatesToClient(stream)
// 	go s.sendLikesUpdatesToClient(stream, func(likeUpdate LikeUpdate) {
// 		var vl string
// 		likeCountKey := fmt.Sprintf("post_likes:%s", likeUpdate.PostID) // Assuming you store likes with this key pattern

// 		vl, _ = s.redisCache.GetKey(stream.Context(), likeCountKey)
// 		log.Println("Val is", vl)

// 		log.Println("Updating in Post stream", likeUpdate.PostID, likeUpdate.UserID)
// 		// Update the metrics of the post with the given postID
// 		// For example, find the post in the cached posts slice and update its like count.
// 		// Then send the updated post to the client.

// 	})

// 	// Wait for the client to disconnect
// 	<-stream.Context().Done()
// 	return stream.Context().Err()
// }

// func convertToPbPost(post *sqlc.GetPostsRecommendationRow) *pb.Post {
// 	var imageUrls []string

// 	err := json.Unmarshal(post.PostImageUrls, &imageUrls)
// 	if err != nil {
// 		// Handle the error
// 		log.Printf("cannot unmarshal data: %s", err)
// 	}

// 	// Convert the sqlc.GetPostsRecommendationRow to the protobuf Post message
// 	// This assumes you have a function to convert from your SQL model to the protobuf model
// 	// You'll need to write this function based on your specific models and protobuf definitions
// 	return &pb.Post{
// 		UserId:       post.UserID.String(),
// 		Username:     post.Username.String,
// 		ProfileImage: post.UserImageUrl.String,
// 		Content:      post.Content,
// 		PostId:       post.ID.String(),
// 		ImageUrls:    imageUrls,
// 		Metrics: &pb.PostMetrics{
// 			Likes:    post.Likes.Int32,
// 			Views:    post.Views.Int32,
// 			Comments: post.Comments.Int32,
// 			Reposts:  post.Reposts.Int32,
// 		},
// 		Timestamp: timestamppb.New(post.CreatedAt.Time),
// 	}
// }

// func (s *SocialMediaServer) sendNewPostsToClient(stream pb.SocialMedia_PostStreamServer, uid uuid.UUID) {
// 	// Subscribe to the Redis channel
// 	pubsub := s.redisCache.Subscribe(stream.Context(), PostChannelKey)
// 	defer pubsub.Close()

// 	for {
// 		msg, err := pubsub.ReceiveMessage(stream.Context())
// 		if err != nil {
// 			if stream.Context().Err() != nil {
// 				return // Stream context was canceled (client disconnected)
// 			}
// 			log.Printf("Error receiving message: %v", err)
// 			continue
// 		}

// 		var post pb.Post
// 		err = json.Unmarshal([]byte(msg.Payload), &post)
// 		if err != nil {
// 			log.Printf("Error unmarshalling post: %v", err)
// 			continue
// 		}

// 		// Check if the current user is following the author of the post
// 		// Note: Assuming `isPostRecommended` checks if the post should be sent to the user.
// 		if isPostRecommended(s.store, stream.Context(), uid) {
// 			if err := stream.Send(&post); err != nil {
// 				log.Printf("Error sending post to stream: %v", err)
// 				return
// 			}
// 		}
// 	}
// }

// func retrieveCachedPosts(cacheClient cache.Cache, ctx context.Context) ([]sqlc.GetPostsRecommendationRow, error) {
// 	cacheKey := PostCacheKey //fmt.Sprintf("posts:%s", uid)
// 	// Assuming cacheClient is a Redis client
// 	cachedPosts, cacheErr := cacheClient.GetKey(ctx, cacheKey) //.Result()
// 	if cacheErr != nil {
// 		return nil, cacheErr
// 	}

// 	var posts []sqlc.GetPostsRecommendationRow
// 	err := json.Unmarshal([]byte(cachedPosts), &posts)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return posts, nil
// }

// func retrievePosts(stream pb.SocialMedia_PostStreamServer, cacheClient cache.Cache, store *sqlc.Store, uid uuid.UUID) ([]sqlc.GetPostsRecommendationRow, error) {

// 	// Assuming a cache is in place, attempt to retrieve posts from cache first
// 	posts, cacheErr := retrieveCachedPosts(cacheClient, stream.Context())
// 	if cacheErr != nil || len(posts) == 0 {
// 		// If cache retrieval fails or cache is empty, fallback to database
// 		var err error

// 		posts, err = store.GetPostsRecommendation(stream.Context(), uid)
// 		if err != nil {
// 			return nil, status.Errorf(codes.Internal, "error fetching posts from db: %s", err)
// 		}
// 		log.Println("Got From SQL", len(posts), " from the db")

// 		// After fetching from the database, add these posts to the cache
// 		cacheErr = addPostsToCache(stream.Context(), cacheClient, posts)
// 		if cacheErr != nil {
// 			log.Printf("failed to add posts to cache: %s", cacheErr)
// 			// Here we just log the error, but you could handle it differently
// 		}
// 	}
// 	return posts, nil
// }

// func addPostsToCache(ctx context.Context, cacheClient cache.Cache, posts []sqlc.GetPostsRecommendationRow) error {
// 	// Convert posts to a cache-friendly format if necessary
// 	// Serialize the posts data for caching (e.g., to JSON)
// 	serializedPosts, err := json.Marshal(posts)
// 	if err != nil {
// 		return err
// 	}

// 	log.Println("Adding post to cache")
// 	// Add the serialized data to the cache with some expiration time
// 	cacheKey := PostCacheKey                                                    //fmt.Sprintf("posts:%s", uid)
// 	cacheErr := cacheClient.SetKey(ctx, cacheKey, serializedPosts, time.Hour*1) //.Result() // Assuming cacheClient is a Redis client
// 	return cacheErr
// }

// func isPostRecommended(store *sqlc.Store, ctx context.Context, uid uuid.UUID) bool {
// 	shouldRecommend, err := store.IsPostRecommended(ctx, uid)
// 	if err != nil {
// 		return false
// 	}
// 	log.Println("Should recommend", shouldRecommend)
// 	return shouldRecommend
// }

// func (s *SocialMediaServer) sendLikesUpdatesToClient(stream pb.SocialMedia_PostStreamServer, updatePostMetrics func(likeUpdate LikeUpdate)) {
// 	likesPubSub := s.redisCache.Subscribe(stream.Context(), LikeChannelKey)
// 	defer likesPubSub.Close()

// 	for {
// 		msg, err := likesPubSub.ReceiveMessage(stream.Context())
// 		if err != nil {
// 			if stream.Context().Err() != nil {
// 				return // Stream context was canceled (client disconnected)
// 			}
// 			log.Printf("Error receiving message: %v", err)
// 			continue
// 		}

// 		var likeUpdate LikeUpdate
// 		err = json.Unmarshal([]byte(msg.Payload), &likeUpdate)
// 		if err != nil {
// 			log.Printf("Error unmarshalling like update: %v", err)
// 			continue
// 		}

// 		updatePostMetrics(likeUpdate)
// 	}
// }

// func (s *SocialMediaServer) updatePostLikes(ctx context.Context, postID string, userID string, increment bool) error {
// 	log.Println("POSTID", postID)
// 	// var err error
// 	var vl string
// 	likeCountKey := fmt.Sprintf("post_likes:%s", postID) // Assuming you store likes with this key pattern

// 	vl, _ = s.redisCache.GetKey(ctx, likeCountKey)
// 	log.Println("Val is", vl)

// 	return nil
// }

// func (s *SocialMediaServer) sendLikesUpdatesToClient1(stream pb.SocialMedia_PostStreamServer) {
// 	likesPubSub := s.redisCache.Subscribe(stream.Context(), LikeChannelKey)
// 	defer likesPubSub.Close()

// 	for {
// 		msg, err := likesPubSub.ReceiveMessage(stream.Context())
// 		if err != nil {
// 			if stream.Context().Err() != nil {
// 				return // Stream context was canceled (client disconnected)
// 			}
// 			log.Printf("Error receiving message: %v", err)
// 			continue
// 		}

// 		var likeUpdate LikeUpdate
// 		err = json.Unmarshal([]byte(msg.Payload), &likeUpdate)
// 		if err != nil {
// 			log.Printf("Error unmarshalling like update: %v", err)
// 			continue
// 		}

// 		// Update post like count in database or cache
// 		err = s.updatePostLikes(stream.Context(), likeUpdate.PostID, likeUpdate.UserID, likeUpdate.Increment)
// 		if err != nil {
// 			// Handle the error appropriately
// 			log.Printf("Error updating post likes: %v", err)
// 		}

// 		// ? Inform the client about the like update if necessary
// 	}
// }

// func (s *SocialMediaServer) sendLikesUpdatesToClientOld(stream pb.SocialMedia_PostStreamServer) {
// 	// Subscribe to the Redis channel for likes
// 	likesPubSub := s.redisCache.Subscribe(stream.Context(), LikeChannelKey)
// 	defer likesPubSub.Close()

// 	for {
// 		msg, err := likesPubSub.ReceiveMessage(stream.Context())
// 		if err != nil {
// 			if stream.Context().Err() != nil {
// 				return // Stream context was canceled (client disconnected)
// 			}
// 			log.Printf("Error receiving message: %v", err)
// 			continue
// 		}

// 		var likeUpdate LikeUpdate // Define LikeUpdate struct according to your needs
// 		err = json.Unmarshal([]byte(msg.Payload), &likeUpdate)
// 		if err != nil {
// 			log.Printf("Error unmarshalling like update: %v", err)
// 			continue
// 		}

// 		// // Send the like update to the client
// 		// if err := stream.Send(&pb.LikeUpdate{
// 		// 	PostId:   likeUpdate.PostId.String(),
// 		// 	NewLikes: likeUpdate.NewLikes,
// 		// }); err != nil {
// 		// 	log.Printf("Error sending likes update to stream: %v", err)
// 		// 	return
// 		// }
// 	}
// }
