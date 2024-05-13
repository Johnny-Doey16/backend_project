package gapi

import (
	"context"
	"database/sql"
	"errors"

	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *SocialMediaServer) CreateRepost(ctx context.Context, req *pb.CreateRepostRequest) (*pb.RepostResponse, error) {
	// claims, ok := stream.Context().Value("payloadKey").(*token.Payload)
	//
	//	if !ok {
	//		return status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	//	}
	uid, _ := services.StrToUUID("ad7a61cd-14c5-4bbe-a3fe-0abdb585898a")

	originalPostID, err := services.StrToUUID(req.OriginalPostId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid original_post_id %v", err)
	}

	// Check post status
	postStatus, err := s.store.CheckPostStatus(ctx, originalPostID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error checking post status: %s", err.Error())
	}

	// Check if the post is active
	if postStatus.PostStatus == "active" {

		// Create transaction
		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "transaction error: %v", err.Error())
		}
		defer tx.Rollback()
		qtx := s.store.WithTx(tx)

		newRepost, total_repost, err := services.RunPostRepostConcurrent(tx, qtx, ctx, originalPostID, uid)
		if err != nil {
			// Handle potential database errors
			return nil, err
		}

		// Add comment count to redis
		err = services.AddRepostCountToRedis(s.redisCache, ctx, RepostsKey, req.GetOriginalPostId(), total_repost)
		if err != nil {
			return nil, err
		}

		// broadcast comment event here
		if err := services.BroadcastRepostEvent(ctx, s.redisCache, req.GetOriginalPostId(), RepostChannelKey, uid); err != nil {
			return nil, err
		}

		// Create the response
		return &pb.RepostResponse{
			Repost: &pb.Repost{
				Id:             int32(newRepost.ID),
				UserId:         newRepost.UserID.String(),
				OriginalPostId: newRepost.OriginalPostID.String(),
				CreatedAt:      timestamppb.New(newRepost.CreatedAt.Time),
			},
		}, nil
	}

	// If the post is not active, return an appropriate response
	return nil, status.Errorf(codes.FailedPrecondition, "Cannot bookmark a %s post", postStatus)
}

func (s *SocialMediaServer) GetRepost(ctx context.Context, req *pb.GetRepostRequest) (*pb.RepostResponse, error) {
	// Convert the id to the appropriate type, e.g., int32 or int64 as required by your database schema.
	repostId := req.Id

	// Call the auto-generated sqlc method to retrieve the repost from the database.
	repost, err := s.store.GetRepostByID(ctx, repostId)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case where the repost is not found.
			return nil, errors.New("repost not found")
		}
		// Handle other potential database errors.
		return nil, err
	}

	// Assuming the `db.GetRepost` method returns a struct of type `db.Repost`, convert it to the protobuf type.
	resp := &pb.RepostResponse{
		Repost: &pb.Repost{
			Id:             int32(repost.ID),
			UserId:         repost.UserID.String(),
			OriginalPostId: repost.OriginalPostID.String(),
			CreatedAt:      timestamppb.New(repost.CreatedAt.Time),
		},
	}

	return resp, nil
}

func (s *SocialMediaServer) GetRepostsByUser(ctx context.Context, req *pb.GetRepostsByUserRequest) (*pb.RepostsResponse, error) {
	// Validate the user_id if needed
	// claims, ok := stream.Context().Value("payloadKey").(*token.Payload)
	//
	//	if !ok {
	//		return status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	//	}
	uid, _ := services.StrToUUID("ad7a61cd-14c5-4bbe-a3fe-0abdb585898a")

	// Retrieve reposts from the database
	reposts, err := s.store.GetRepostsByUserID(ctx, uid)
	if err != nil {
		// Handle potential database errors
		return nil, err
	}

	// Convert the retrieved reposts to protobuf message format
	var pbReposts []*pb.Repost
	for _, repost := range reposts {
		pbReposts = append(pbReposts, &pb.Repost{
			Id:             int32(repost.ID),
			UserId:         repost.UserID.String(),
			OriginalPostId: repost.OriginalPostID.String(),
			CreatedAt:      timestamppb.New(repost.CreatedAt.Time),
		})
	}

	// Create the response
	resp := &pb.RepostsResponse{
		Reposts: pbReposts,
	}

	return resp, nil
}

func (s *SocialMediaServer) DeleteRepost(ctx context.Context, req *pb.DeleteRepostRequest) (*pb.DeleteRepostResponse, error) {
	// Convert the id to the appropriate type, e.g., int32 or int64 as required by your database schema.
	repostId := req.Id

	// Call the auto-generated sqlc method to delete the repost from the database.
	err := s.store.DeleteRepost(ctx, repostId)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case where the repost is not found.
			return nil, errors.New("repost not found")
		}
		// Handle other potential database errors.
		return nil, err
	}

	// Construct the response indicating successful deletion.
	resp := &pb.DeleteRepostResponse{
		Success: true,
	}

	return resp, nil
}
