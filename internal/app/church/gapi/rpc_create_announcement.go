package gapi

import (
	"context"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/constants"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ChurchServer) CreateAnnouncement(ctx context.Context, req *pb.CreateAnnouncementRequest) (*pb.CreateAnnouncementResponse, error) {
	// TODO: Sanitize inputs

	// TODO: Check whether to Add constraints in db so it is only churches that can make announcements

	// claims, ok := ctx.Value("payloadKey").(*token.Payload)
	// if !ok {
	// 	return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	// }

	uid, _ := services.StrToUUID("1f1adf60-988a-46c8-83a2-9c30d5c146bd")

	newsId, err := uuid.NewRandom()
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "an unexpected error occurred")
	}

	err = s.store.CreateAnnouncements(ctx, sqlc.CreateAnnouncementsParams{
		ID:      newsId,
		UserID:  uid,
		Title:   req.GetTitle(),
		Content: req.GetContent(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Error while creating announcement %v", err)
	}

	membersIds, err := s.store.GetChurchMembersUid(ctx, uid)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Error while creating announcement %v", err)
	}

	services.SendNotification(s.taskDistributor, ctx, uid, newsId, membersIds, constants.NotificationChurchAnnouncement, "Announcement", "", req.GetContent(), "")

	return &pb.CreateAnnouncementResponse{
		NewsId: newsId.String(),
		Msg:    "Announcement successfully made.",
	}, nil
}
