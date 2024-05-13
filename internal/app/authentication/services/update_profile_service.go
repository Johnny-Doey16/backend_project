package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/pb"
)

func UpdateProfile(ctx context.Context, store *sqlc.Store, uid uuid.UUID, req *pb.UpdateProfileRequest) error {
	// Sanitize inputs
	err := sanitizeInputs(req.GetImageUrl(), req.GetFirstName(), req.GetLastName())
	if err != nil {
		return fmt.Errorf("error updating profile %s", err)
	}

	// Send data to db
	// err = store.UpdateUserProfile(ctx, sqlc.UpdateUserProfileParams{
	// 	UserID:    uid,
	// 	FirstName: sql.NullString{String: req.GetFirstName(), Valid: req.FirstName != nil},
	// 	LastName:  sql.NullString{String: req.GetLastName(), Valid: req.LastName != nil},
	// 	ImageUrl:  sql.NullString{String: req.GetImageUrl(), Valid: req.ImageUrl != nil},
	// 	UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	// })
	// if err != nil {
	// 	return fmt.Errorf("error updating profile %s", err)
	// }
	return nil
}

func sanitizeInputs(imageUrl, firstName, lastName string) error {
	// TODO: Sanitize inputs that are not nil.
	return nil
}
