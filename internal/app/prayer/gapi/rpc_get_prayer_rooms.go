package gapi

import (
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/prayer/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *PrayerServer) GetUserPrayerRooms(_ *pb.GetUserRoomsRequest, stream pb.PrayerService_GetUserPrayerRoomsServer) error {

	// TODO: Paginate
	claims, ok := stream.Context().Value("payloadKey").(*token.Payload)
	if !ok {
		return status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	rooms, err := s.store.GetPrayerRooms(stream.Context(), sqlc.GetPrayerRoomsParams{
		UserID: claims.UserId,
		Offset: 0,
		Limit:  30, // todo: add to proto
	})
	if err != nil {
		return status.Errorf(codes.Internal, "error fetching prayer rooms %v", err)
	}

	for i := 0; i < len(rooms); i++ {

		stream.Send(&pb.PrayerRoom{
			RoomId:    rooms[i].RoomID,
			Name:      rooms[i].Name.String,
			CreatedAt: timestamppb.New(rooms[i].CreatedAt.Time),
			StartTime: timestamppb.New(rooms[i].StartTime.Time),
			EndTime:   timestamppb.New(rooms[i].EndTime.Time),
			IsAuthor:  rooms[i].AuthorID == claims.UserId,
			Status:    string(rooms[i].ParticipantStatus.ParticipantStatus),
		})
	}

	return nil
}
