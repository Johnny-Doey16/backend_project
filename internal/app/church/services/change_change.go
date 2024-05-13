package services

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UpdateOrCreateChurchMembership(ctx context.Context, store *sqlc.Store, userID uuid.UUID, churchID, membershipID int32, hasMembership bool) error {
	var err error
	if !hasMembership || membershipID == 0 {
		log.Println("Creating new membership")
		_, err = store.CreateChurchForUser(ctx, sqlc.CreateChurchForUserParams{
			UserID:   userID,
			ChurchID: int32(churchID),
		})
	} else {
		log.Println("Updating old membership...")
		_, err = store.UpdateChurchForUser(ctx, sqlc.UpdateChurchForUserParams{
			ChurchID: int32(churchID),
			UserID:   userID,
		})
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return status.Errorf(codes.NotFound, "you can only join a church that is same denomination as your denomination. consider changing your denomination or join a different church")
		}

		return status.Errorf(codes.Internal, "error updating/creating church membership: %s", err)
	}
	return nil
}
