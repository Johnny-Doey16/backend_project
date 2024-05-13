package gapi

import (
	"context"
	"database/sql"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/prayer/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *PrayerServer) UpdatePrayerRoom(ctx context.Context, req *pb.PrayerRoom) (*pb.PrayerRoomResponse, error) {
	// TODO: in the sql query let it be the author of the prayer that can edit it

	var startedAtTime time.Time
	var endAtTime time.Time
	var err error
	if req.CreatedAt != nil {
		startedAtTime, err = ptypes.Timestamp(req.CreatedAt)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error converting time %v", err)
		}
	}

	if req.EndTime != nil {
		endAtTime, err = ptypes.Timestamp(req.CreatedAt)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error converting time %v", err)
		}
	}

	err = s.store.UpdatePrayerRoom(ctx, sqlc.UpdatePrayerRoomParams{
		RoomID:    req.RoomId,
		Name:      sql.NullString{Valid: req.Name != "", String: req.Name},
		StartTime: sql.NullTime{Valid: req.CreatedAt != nil, Time: startedAtTime},
		EndTime:   sql.NullTime{Valid: req.EndTime != nil, Time: endAtTime},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error updating room %v", err)
	}

	// TODO: Update in the signalling server

	return &pb.PrayerRoomResponse{
		RoomId:   req.RoomId,
		RoomName: req.Name,
		Msg:      "Room updated successfully",
	}, nil
}
