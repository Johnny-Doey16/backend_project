package gapi2

/*
func (s *SocialMediaServer) PostStream(stream pb.SocialMedia_PostStreamServer) error {
	// _, ok := stream.Context().Value("payloadKey").(*token.Payload)
	//
	//	if !ok {
	//		return status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	//	}

	uid, _ := services.StrToUUID("ad7a61cd-14c5-4bbe-a3fe-0abdb585898a")

	batchSize := 40
	offset := 0

	// Create a channel for this client to receive posts
	postChannel := make(chan pb.Post)
	defer close(postChannel)

	// Subscribe the client to receive posts
	s.subscribe <- postChannel

	posts, err := retrievePosts(stream, s.redisCache, s.store, uid)
	if err != nil {
		return err
	}

	log.Println("Total is: ", len(posts))

	var imageUrls []string
	for _, v := range posts {
		post := v

		err = json.Unmarshal(post.PostImageUrls, &imageUrls)
		if err != nil {
			// Handle the error
			log.Printf("cannot unmarshal data: %s", err)
		}
		// log.Println("Image urls", imageUrls[0])

		if err := stream.Send(
			&pb.Post{
				PostId:       post.ID.String(),
				Content:      post.Content,
				FirstName:    post.FirstName.String,
				Username:     post.Username.String,
				ProfileImage: post.UserImageUrl.String,
				UserId:       post.UserID.String(),
				Timestamp:    timestamppb.New(post.CreatedAt.Time),
				ImageUrls:    imageUrls,
			},
		); err != nil {
			return status.Errorf(codes.Internal, "error fetching posts: %s", err)
		}
	}
	offset += batchSize

	// Start a goroutine to send new posts to the client
	go s.sendNewPostsToClient(stream, postChannel, uid)

	// Wait for the client to disconnect
	<-stream.Context().Done()

	return stream.Context().Err()
}

func (s *SocialMediaServer) sendNewPostsToClient(stream pb.SocialMedia_PostStreamServer, postChannel chan pb.Post, uid uuid.UUID) {
	for {
		select {
		case post, ok := <-postChannel:
			if !ok {
				// The client channel is closed; unsubscribe the client
				s.unsubscribe <- postChannel
				return
			}
			// Check if the current user is following the author of the post
			if isPostRecommended(s.store, stream.Context(), uid) { //isFollowing(s.store, stream.Context(), uid, post.UserId) {
				if err := stream.Send(&post); err != nil {
					// Handle the error if necessary.
					return
				}
			}
		case <-stream.Context().Done():
			// Client disconnected; unsubscribe the client
			s.unsubscribe <- postChannel
			return
		}
	}
}

func retrieveCachedPosts(cacheClient cache.Cache, ctx context.Context, uid uuid.UUID) ([]sqlc.GetPostsRecommendationRow, error) {
	cacheKey := fmt.Sprintf("posts:%s", uid)
	// Assuming cacheClient is a Redis client
	cachedPosts, cacheErr := cacheClient.GetKey(ctx, cacheKey) //.Result()
	if cacheErr != nil {
		return nil, cacheErr
	}

	var posts []sqlc.GetPostsRecommendationRow
	err := json.Unmarshal([]byte(cachedPosts), &posts)
	if err != nil {
		return nil, err
	}
	log.Println("Got", len(posts), " from the db")
	// log.Println("Got this from cache", posts)

	return posts, nil
}

func retrievePosts(stream pb.SocialMedia_PostStreamServer, cacheClient cache.Cache, store *sqlc.Store, uid uuid.UUID) ([]sqlc.GetPostsRecommendationRow, error) {

	// posts, err := s.store.GetPosts(stream.Context(), sqlc.GetPostsParams{
	// 	Limit:  int32(batchSize),
	// 	Offset: int32(offset),
	// })

	// Assuming a cache is in place, attempt to retrieve posts from cache first
	posts, cacheErr := retrieveCachedPosts(cacheClient, stream.Context(), uid)
	if cacheErr != nil || len(posts) == 0 {
		// If cache retrieval fails or cache is empty, fallback to database
		var err error

		// 	posts, err := s.store.GetPostsByFollowing(stream.Context(), uid)
		// posts, err := s.store.GetPostsRecommendationOLD(stream.Context(), uid)
		posts, err = store.GetPostsRecommendation(stream.Context(), uid)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error fetching posts from db: %s", err)
		}
		log.Println("Got From SQL", len(posts), " from the db")

		// After fetching from the database, add these posts to the cache
		cacheErr = addPostsToCache(stream.Context(), cacheClient, uid, posts)
		if cacheErr != nil {
			log.Printf("failed to add posts to cache: %s", cacheErr)
			// Here we just log the error, but you could handle it differently
		}
	}
	return posts, nil
}

func addPostsToCache(ctx context.Context, cacheClient cache.Cache, uid uuid.UUID, posts []sqlc.GetPostsRecommendationRow) error {
	// Convert posts to a cache-friendly format if necessary
	// Serialize the posts data for caching (e.g., to JSON)
	serializedPosts, err := json.Marshal(posts)
	if err != nil {
		return err
	}

	log.Println("Adding post to cache")
	// Add the serialized data to the cache with some expiration time
	cacheKey := fmt.Sprintf("posts:%s", uid)
	cacheErr := cacheClient.SetKey(ctx, cacheKey, serializedPosts, time.Hour*1) //.Result() // Assuming cacheClient is a Redis client
	return cacheErr
}

// ? Fix to allow author also get this stream
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

func isPostRecommended(store *sqlc.Store, ctx context.Context, uid uuid.UUID) bool {
	shouldRecommend, err := store.IsPostRecommended(ctx, uid)
	if err != nil {
		return false
	}
	log.Println("Should recommend", shouldRecommend)
	return shouldRecommend
}
*/
