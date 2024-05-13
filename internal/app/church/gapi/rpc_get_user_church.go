package gapi

import (
	"context"
	"database/sql"

	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ChurchServer) GetUserChurch(ctx context.Context, _ *pb.GetUserChurchRequest) (*pb.GetUserChurchResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	church, err := s.store.GetUserChurch(ctx, claims.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "you are not register in a church")
		}
		return nil, status.Errorf(codes.Internal, "unable to retrieve Church %v", err.Error())
	}
	return &pb.GetUserChurchResponse{
		Church: &pb.Church{
			AuthId:          church.ChurchAuthID.String(),
			Id:              church.ChurchID,
			Name:            church.ChurchName,
			ImageUrl:        church.ImageUrl.String,
			Username:        church.Username.String,
			DenominationId:  church.DenominationID,
			Email:           church.Email,
			Phone:           church.Phone.String,
			FollowingCount:  int64(church.FollowingCount.Int32),
			FollowerCount:   int64(church.FollowersCount.Int32),
			MembershipCount: int64(church.MembersCount.Int32),
			IsVerified:      church.IsVerified.Bool,
			PostCount:       int64(church.PostsCount.Int32),
			About:           church.About.String,
			Website:         church.Website.String,
			HeaderImageUrl:  church.HeaderImageUrl.String,
			AccountName:     church.AccountName.String,
			AccountNumber:   church.AccountNumber.String,
			BankName:        church.BankName.String,
			Location: &pb.Location{
				Country:    church.Country,
				State:      church.State,
				City:       church.City,
				Address:    church.Address,
				PostalCode: church.Postalcode,
				Lga:        church.Lga,
			},
		},
		IsFollowing: church.IsFollowing,
		IsFollowed:  church.IsFollowed,
	}, nil
}
