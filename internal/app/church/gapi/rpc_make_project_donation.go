package gapi

import (
	"context"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO: Fix data input type issue
func (s *ChurchServer) CreateProjectDonate(ctx context.Context, req *pb.ProjectDonationRequest) (*pb.ProjectDonateResponse, error) {
	// claims, ok := ctx.Value("payloadKey").(*token.Payload)

	// if !ok {
	// 	return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	// }
	userId, err := services.StrToUUID("558c60aa-977f-4d38-885b-e813656371ac")
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "parsing uid: %s", err)
	}

	err = s.store.CreateProjectsDonation(ctx, sqlc.CreateProjectsDonationParams{
		ProjectID: int32(req.GetProjectId()),
		UserID:    userId,
		// DonationAmount: req.GetAmount(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Donation error: %s", err)
	}

	return &pb.ProjectDonateResponse{
		Msg: "success adding donation",
	}, nil
}
