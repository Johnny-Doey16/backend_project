package gapi2

/**
func (s *SocialMediaServer) broadcast(post *pb.Post) {
	s.subscribers.Range(func(key, value interface{}) bool {
		subscriber, ok := key.(chan pb.Post)
		if !ok {
			return true // Continue to the next subscriber if type assertion fails
		}

		select {
		case subscriber <- *post:
			// Message sent successfully
		default:
			// If the client is slow, we skip it, but consider removing the subscriber
			// from the map to prevent future failed send attempts.
		}

		// Check if the channel is closed and if so, delete it from the map.
		// You can do this by trying to receive from the channel with a default case.
		// If the receive succeeds or the default case is hit, then the channel is not closed.
		// If the receive blocks, then a separate goroutine is used to check if the channel is closed.
		select {
		case _, ok := <-subscriber:
			if !ok {
				s.subscribers.Delete(key) // Delete the subscriber since the channel is closed
			}
		default:
			// Channel is not ready for receiving which means it's not closed.
		}

		return true // Continue to the next subscriber
	})
}

func (s *SocialMediaServer) oldBroadcast(post *pb.Post) {
	s.subscribers.Range(func(key, _ interface{}) bool { //value interface{}) bool {
		subscriber, ok := key.(chan pb.Post)
		if ok {
			select {
			case subscriber <- *post:
			default:
				// Skip if the client is slow to consume messages
			}
		}
		return true
	})
}

// Check if image and texts are safe to upload
func (s *SocialMediaServer) CreatePost(stream pb.SocialMedia_CreatePostServer) error {
	ctx, cancel := context.WithTimeout(stream.Context(), time.Second*10) // Adjust the timeout as needed
	defer cancel()

	// Check whether to remove or leave
	s.mu.Lock()
	defer s.mu.Unlock()

	// claims, ok := ctx.Value("payloadKey").(*token.Payload)
	// if !ok {
	// 	return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	// }
	userId, err := services.StrToUUID("9cd93ffa-5835-414f-95c7-f5e94608c71c")
	if err != nil {
		return status.Errorf(codes.Internal, "parsing uid: %s", err)
	}

	imageUrls, captions, content, err := services.GetCreatePostStreamData(stream, s.imageStore)
	if err != nil {
		return err
	}

	postID, sqlPost, err := services.CreatePostInDB(ctx, s.db, s.store, userId, content, imageUrls, captions)
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

	if err := s.cachePost(ctx, &post); err != nil {
		// Handle the error but you might not want to fail the whole operation
		// because the main action (inserting into DB) was successful.
		// Log this error instead of returning it.
		log.Printf("Failed to cache the post: %v\n", err)
	}

	// Broadcast the post to subscribers
	s.broadcast(&post)
}

func (s *SocialMediaServer) cachePost(ctx context.Context, post *pb.Post) error {
	// Use the post ID as the key and serialize the post object as the value.
	// You can use JSON or any other serialization method of your choice.
	// Here's an example using JSON:
	serializedPost, err := json.Marshal(post)
	if err != nil {
		return err
	}

	log.Println("Adding post to cache")
	// Cache the serialized post in Redis with an expiration time (e.g., 24 hours).
	// The key pattern here is "post:{postID}".
	return s.redisCache.SetKey(ctx, fmt.Sprintf("posts:%s", post.PostId), serializedPost, 24*time.Hour) //.Err()
}

*/
