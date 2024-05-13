package gapi

import (
	"context"

	"github.com/steve-mir/diivix_backend/internal/app/prayer/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *PrayerServer) DeletePrayerRoom(ctx context.Context, req *pb.PrayerRoomId) (*pb.PrayerRoomResponse, error) {
	// TODO: in the sql query let it be the author of the prayer that can delete it
	err := s.store.DeletePrayerRoom(ctx, req.RoomId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error deleting prayer room %v", err)
	}

	// TODO: Delete room in signalling server

	return &pb.PrayerRoomResponse{
		Msg:    "Prayer room deleted successfully",
		RoomId: req.RoomId,
	}, nil
}
