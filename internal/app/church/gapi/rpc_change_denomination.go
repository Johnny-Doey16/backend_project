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

func (s *ChurchServer) ChangeDenominationMembership(ctx context.Context, req *pb.MembershipChangeRequest) (*pb.MembershipChangeResponse, error) {
	// Step 1: Get user ID from the request
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	// Step 2: Check if denomination_id is provided in the request
	denominationID := req.GetDenominationId()
	if denominationID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "denomination ID is required")
	}

	membership, err := s.store.GetUserAndMembership(ctx, claims.UserId)
	if err != nil && err != sql.ErrNoRows {
		return nil, status.Errorf(codes.Internal, "error checking current denomination membership: %s", err)
	}

	hasMembership := membership != (sqlc.GetUserAndMembershipRow{})
	if hasMembership {
		if membership.DenominationID.Int32 == denominationID {
			return nil, status.Errorf(codes.AlreadyExists, "user is already a member of this denomination %+v", membership.Active.Bool)
		}

		// Enforce the rule that denomination can only be changed once a year
		if membership.Active.Bool && time.Since(membership.JoinDate.Time) < (time.Hour*24*365) {
			return nil, status.Errorf(codes.FailedPrecondition, "denomination can only be changed once a year")
		}
	}

	// Update or create the membership
	if err := services.UpdateOrCreateDenominationMembership(ctx, s.store, claims.UserId, denominationID, membership.DenominationID_2.Int32, hasMembership); err != nil {
		return nil, err
	}

	// Successfully updated or created membership
	return &pb.MembershipChangeResponse{}, nil
}
