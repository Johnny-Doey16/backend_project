package gapi

import (
	"context"
	"database/sql"
	"time"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"github.com/steve-mir/diivix_backend/internal/app/church/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ChurchServer) ChangeChurchMembership(ctx context.Context, req *pb.MembershipChangeRequest) (*pb.MembershipChangeResponse, error) {

	// Step 1: Get user ID from the request
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	// Step 2: Check if denomination_id is provided in the request
	churchID := req.GetChurchId()
	if churchID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "church ID is required")
	}

	// membership, err := s.store.GetUserAndChurchMembership(ctx, claims.UserId)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return nil, status.Errorf(codes.NotFound, "membership not found")
	// 	}
	// 	return nil, status.Errorf(codes.Internal, "error checking current church membership: %s", err)
	// }

	membership, err := s.store.GetUserAndChurchMembership(ctx, claims.UserId)
	if err != nil && err != sql.ErrNoRows {
		return nil, status.Errorf(codes.Internal, "error checking current church membership: %s", err)
	}

	hasMembership := membership != (sqlc.GetUserAndChurchMembershipRow{})
	if hasMembership {
		if membership.ChurchID.Int32 == churchID {
			return nil, status.Errorf(codes.AlreadyExists, "user is already a member of this church %+v", membership.Active.Bool)
		}

		// Enforce the rule that denomination can only be changed once a year
		if membership.Active.Bool && time.Since(membership.JoinDate.Time) < (time.Hour*24*182) {
			return nil, status.Errorf(codes.FailedPrecondition, "church can only be changed once in 6 months")
		}
	}

	// ! Update or create the membership
	// if !hasMembership || membership.ChurchID.Int32 == 0 {
	// 	log.Println("Creating new membership")
	// 	_, err = s.store.CreateChurchForUser(ctx, sqlc.CreateChurchForUserParams{
	// 		UserID:   userID,
	// 		ChurchID: int32(churchID),
	// 	})
	// } else {
	// 	log.Println("Updating old membership...")
	// 	_, err = s.store.UpdateChurchForUser(ctx, sqlc.UpdateChurchForUserParams{
	// 		ChurchID: int32(churchID),
	// 		UserID:   userID,
	// 	})
	// }

	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return nil, status.Errorf(codes.NotFound, "you can only join a church that is same denomination as your denomination. consider changing your denomination or join a different church")
	// 	}

	// 	return nil, status.Errorf(codes.Internal, "error updating/creating church membership: %s", err)
	// }

	// Successfully updated or created membership
	if err := services.UpdateOrCreateChurchMembership(ctx, s.store, claims.UserId, churchID, membership.ChurchID.Int32, hasMembership); err != nil {
		return nil, err
	}
	return &pb.MembershipChangeResponse{}, nil
}
