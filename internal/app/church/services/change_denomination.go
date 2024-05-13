package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UpdateOrCreateDenominationMembership(ctx context.Context, store *sqlc.Store, userID uuid.UUID, denominationID, membershipID int32, hasMembership bool) error {
	var err error
	if !hasMembership || membershipID == 0 {
		err = store.CreateDenominationForUser(ctx, sqlc.CreateDenominationForUserParams{
			UserID:         userID,
			DenominationID: int32(denominationID),
		})
	} else {
		err = store.UpdateDenominationForUser(ctx, sqlc.UpdateDenominationForUserParams{
			DenominationID: int32(denominationID),
			UserID:         userID,
		})
	}

	if err != nil {
		return status.Errorf(codes.Internal, "error updating/creating denomination membership: %s", err)
	}
	return nil
}
