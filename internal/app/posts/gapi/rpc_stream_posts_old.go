package gapi

// func (s *SocialMediaServer) PostStream(stream pb.SocialMedia_PostStreamServer) error {
// 	// claims, ok := stream.Context().Value("payloadKey").(*token.Payload)

// 	// if !ok {
// 	// 	return status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
// 	// }
// 	userId, err := services.StrToUUID("558c60aa-977f-4d38-885b-e813656371ac") // ad7a61cd-14c5-4bbe-a3fe-0abdb585898a
// 	if err != nil {
// 		return status.Errorf(codes.Internal, "parsing uid: %s", err)
// 	}

// 	// Ensure that the page number and page size are positive.
// 	pageNumber := 1
// 	pageSize := 2

// 	// Calculate the offset based on the page number and page size.
// 	offset := (pageNumber - 1) * pageSize
// 	log.Println("OFFSET", offset)

// 	// Retrieve posts either from cache or database
// 	posts, err := services.RetrievePosts(stream.Context(), s.store, s.redisCache, userId, PostCacheKey, PostCacheExpr, int32(offset), int32(pageSize))
// 	if err != nil {
// 		return err
// 	}

// 	// Send the retrieved posts to the client first
// 	if err := services.SendPostToClient(posts, s.redisCache, stream, userId, PostLikesKey); err != nil {
// 		return status.Errorf(codes.Internal, "error sending post: %s", err)
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
// 	go services.ListenForRefresh(stream, s.store, s.redisCache, userId, PostChannelKey, PostLikesKey)

// 	go services.SendNewPostsToClientOld(stream, s.store, s.redisCache, userId, PostChannelKey, PostLikesKey)
// 	// go services.SendNewPostsToClient2(stream, s.store, s.redisCache, userId, PostChannelKey, PostLikesKey)
// 	// go services.SendNewPostsToClient(stream, s.store, s.redisCache, userId, PostChannelKey, PostLikesKey)

// 	go services.SendLikesUpdatesToClient(s.redisCache, stream, LikeChannelKey, PostLikesKey)

// 	go services.SendCommentsUpdatesToClient(s.redisCache, stream, CommentChannelKey, PostCommentsKey)

// 	go services.SendRepostsUpdatesToClient(s.redisCache, stream, RepostChannelKey, RepostsKey)

// 	go services.SendViewsUpdatesToClient(s.redisCache, stream, ViewChannelKey)

// 	// Wait for the client to disconnect
// 	<-stream.Context().Done()
// 	return stream.Context().Err()
// }
