package gapi

import (
	"log"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ChurchServer) GetAnnouncementsForUser(req *pb.GetAnnouncementsForUserRequest, stream pb.ChurchService_GetAnnouncementsForUserServer) error {
	// claims, ok := ctx.Value("payloadKey").(*token.Payload)
	// if !ok {
	// 	return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	// }

	uid, _ := services.StrToUUID("558c60aa-977f-4d38-885b-e813656371ac")

	// Calculate the offset based on the page number and page size
	offset := (req.GetPageNumber() - 1) * req.GetPageSize()
	limit := req.GetPageSize()

	// Calculate the offset based on the page number and page size.
	if req.PageNumber < 1 {
		req.PageNumber = 1
	}
	if limit < 1 {
		limit = 10 // Default page size to 10 if not specified or negative
	}

	news, err := s.store.GetAnnouncementsForUser(stream.Context(), sqlc.GetAnnouncementsForUserParams{
		UserID: uid,
		Limit:  limit,
		Offset: int32(offset),
	})
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	log.Println("NEWS", news)

	var hasMore bool = true
	if len(news) < int(limit) {
		hasMore = false
	}

	for _, post := range news {
		stream.Send(&pb.GetAnnouncementsResponse{
			Post: &pb.Announcement{
				PostId:       post.ID.String(),
				UserId:       post.UserID.String(),
				Title:        post.Title,
				Content:      post.Content,
				Name:         post.ChurchName,
				Username:     post.Username.String,
				ProfileImage: post.ChurchImageUrl.String,
				IsVerified:   post.ChurchVerified.Bool,
				Timestamp:    timestamppb.New(post.CreatedAt.Time),
			},
			HasMore: hasMore,
		})
	}

	return nil
}
