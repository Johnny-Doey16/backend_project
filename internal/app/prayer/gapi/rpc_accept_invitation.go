package gapi

import (
	"context"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/prayer/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	accepted = "accepted"
	decline  = "declined"
)

func (s *PrayerServer) AcceptInvitation(ctx context.Context, req *pb.PrayerRoomId) (*pb.InvitationResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	err := s.store.UpdatePrayerInvitation(ctx, sqlc.UpdatePrayerInvitationParams{
		RoomID: req.RoomId,
		UserID: claims.UserId,
		Status: sqlc.NullParticipantStatus{ParticipantStatus: accepted, Valid: true},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error accepting invitation %v", err)
	}

	return &pb.InvitationResponse{
		Msg:      "invitation accepted",
		Accepted: true,
	}, nil
}
