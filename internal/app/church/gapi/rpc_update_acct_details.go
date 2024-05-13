package gapi

import (
	"context"
	"database/sql"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ChurchServer) UpdateChurchAccountDetails(ctx context.Context, req *pb.UpdateAccountDetailsRequest) (*pb.UpdateAccountDetailsResponse, error) {
	// ! Check if account name an church name are both similar
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	err := s.store.UpdateUserAccountDetails(ctx, sqlc.UpdateUserAccountDetailsParams{
		UserID:        claims.UserId,
		AccountName:   sql.NullString{String: req.GetAccountNumber(), Valid: true},
		AccountNumber: sql.NullString{String: req.GetAccountNumber(), Valid: true},
		BankName:      sql.NullString{String: req.GetBankName(), Valid: true},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error updating account details %s", err.Error())
	}

	return &pb.UpdateAccountDetailsResponse{
		Msg: "account details updated successfully",
	}, nil
}
