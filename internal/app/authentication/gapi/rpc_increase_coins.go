package gapi

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/pb"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) IncreaseTotalCoin(ctx context.Context, req *pb.IncreaseTotalCoinRequest) (*pb.IncreaseTotalCoinResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	total, err := increaseCoin(ctx, s.db, claims.UserId, req.GetAmountToIncrease())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error subtracting %v", err)
	}

	return &pb.IncreaseTotalCoinResponse{
		NewTotalCoin: total,
	}, nil
}

func increaseCoin(ctx context.Context, db *sql.DB, uid uuid.UUID, amount float64) (float64, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, status.Errorf(codes.Internal, "transaction error: %v", err.Error())
	}
	defer tx.Rollback()

	// Executing the SQL code within the transaction
	query := fmt.Sprintf(`
        UPDATE accounts
        SET total_coin = total_coin + %f
        WHERE user_id = '%s'
        RETURNING total_coin;
    `, amount, uid)

	var newTotal float64
	err = tx.QueryRow(query).Scan(&newTotal)
	if err != nil {
		// If there's an error, rollback the transaction
		tx.Rollback()
		return 0, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return newTotal, nil
}
