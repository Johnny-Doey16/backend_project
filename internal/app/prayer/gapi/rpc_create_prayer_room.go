package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/constants"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	post_service "github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"github.com/steve-mir/diivix_backend/internal/app/prayer/pb"
	"github.com/steve-mir/diivix_backend/internal/app/prayer/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO: Refactor
func (s *PrayerServer) CreatePrayerRoom(ctx context.Context, req *pb.CreatePrayerRoomRequest) (*pb.PrayerRoomResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	// Verify start time is a future time
	startedAtTime, err := ptypes.Timestamp(req.GetPrayerRoom().GetStartTime())
	if err != nil || startedAtTime.Before(time.Now()) {
		return nil, status.Errorf(codes.InvalidArgument, "start time should be a future time passed: %+v. Got: %v", req.GetPrayerRoom().GetStartTime(), startedAtTime)
		// return nil, status.Error(codes.InvalidArgument, "start time should be a future time")
	}

	// Verify end time is at least 45 minutes more than start time and not greater than 8 hours
	endAtTime, err := ptypes.Timestamp(req.GetPrayerRoom().GetEndTime())
	if err != nil || endAtTime.Before(startedAtTime.Add(45*time.Minute)) || endAtTime.After(startedAtTime.Add(8*time.Hour)) {
		return nil, status.Error(codes.InvalidArgument, "end time should be at least 45 minutes more than start time and not greater than 8 hours")
	}

	// roomId, err := services.GenerateRoomId()
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "error creating room %v", err)
	// }

	roomId, err := services.CreateRoomInSignallingServer("", req.GetPrayerRoom().GetName())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error calling signalling server %v", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transaction error: %v", err.Error())
	}
	defer tx.Rollback()
	qtx := s.store.WithTx(tx)

	// Create prayer room
	err = qtx.CreatePrayerRoom(ctx, sqlc.CreatePrayerRoomParams{
		RoomID:    roomId,
		AuthorID:  claims.UserId,
		Name:      sql.NullString{String: req.GetPrayerRoom().GetName(), Valid: true},
		StartTime: sql.NullTime{Valid: true, Time: startedAtTime},
		EndTime:   sql.NullTime{Valid: true, Time: endAtTime},
	})
	if err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "error creating prayer room %v %s", err, roomId)
	}

	// Get users' ID from their usernames
	participantsUIDs, err := post_service.GetUserIDsFromUsernames(tx, claims.Username, req.ParticipantUsernames)
	if err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "error getting uids %v", err)
	}

	err = qtx.CreatePrayerParticipants(ctx, sqlc.CreatePrayerParticipantsParams{
		RoomID:  roomId,
		Column1: participantsUIDs,
	})
	if err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "error creating prayer participants %v", err)
	}

	// send notification to all users
	if len(participantsUIDs) > 0 {
		notificationMessage := fmt.Sprintf("%s Invited you to %s prayer room. Starting %+v", claims.Username, req.GetPrayerRoom().GetName(), req.GetPrayerRoom().StartTime)
		post_service.SendNotification(s.taskDistributor, ctx, claims.UserId, uuid.Nil, participantsUIDs, constants.NotificationPrayerInvite, "Prayer", "", notificationMessage, roomId)
	}

	return &pb.PrayerRoomResponse{
		Msg:      "Prayer room created created",
		RoomId:   roomId,
		RoomName: req.GetPrayerRoom().GetName(),
	}, tx.Commit()
}
