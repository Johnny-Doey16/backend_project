package gapi

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/constants"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	post_service "github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"github.com/steve-mir/diivix_backend/internal/app/prayer/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO: Reconcile similar functions in create prayer
func (s *PrayerServer) InviteParticipant(ctx context.Context, req *pb.InviteParticipantRequest) (*pb.PrayerRoomResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transaction error: %v", err.Error())
	}
	defer tx.Rollback()
	qtx := s.store.WithTx(tx)

	// Get users' ID from their usernames
	participantsUIDs, err := post_service.GetUserIDsFromUsernames(tx, claims.Username, req.ParticipantUsernames)
	if err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "error getting uids %v", err)
	}

	err = qtx.CreatePrayerParticipants(ctx, sqlc.CreatePrayerParticipantsParams{
		RoomID:  req.RoomId,
		Column1: participantsUIDs,
	})
	if err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "error creating prayer participants %v", err)
	}

	// send notification to all users
	if len(participantsUIDs) > 0 {
		notificationMessage := fmt.Sprintf("%s Invited you to %s prayer room.", claims.Username, req.RoomName)
		post_service.SendNotification(s.taskDistributor, ctx, claims.UserId, uuid.Nil, participantsUIDs, constants.NotificationPrayerInvite, "Prayer", "", notificationMessage, req.RoomId)
	}

	return &pb.PrayerRoomResponse{
		Msg:      "New invite sent",
		RoomId:   req.RoomId,
		RoomName: req.RoomName,
	}, tx.Commit()
}
