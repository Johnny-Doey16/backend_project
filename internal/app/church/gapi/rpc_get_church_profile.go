package gapi

import (
	"context"
	"log"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	post_service "github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ChurchServer) GetChurchProfile(ctx context.Context, req *pb.GetChurchProfileRequest) (*pb.GetChurchProfileResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}
	churchAuthID, _ := post_service.StrToUUID(req.GetAuthId())

	church, err := s.store.GetChurchProfile(ctx, sqlc.GetChurchProfileParams{
		AuthID:         churchAuthID,
		FollowerUserID: claims.UserId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to retrieve profile")
	}
	log.Println("Is following", church.IsFollowing, "Is followed by church", church.IsFollowed, "Is church member", church.IsMember)
	return &pb.GetChurchProfileResponse{
		IsFollowing: church.IsFollowing,
		IsFollowed:  church.IsFollowed,
		IsMember:    church.IsMember,
		Church: &pb.Church{
			AuthId:          church.AuthID.String(),
			Id:              church.ID,
			Name:            church.Name,
			ImageUrl:        church.ImageUrl.String,
			Username:        church.Username.String,
			DenominationId:  church.DenominationID,
			MembershipCount: int64(church.MembersCount.Int32),
			Phone:           church.Phone.String,
			IsVerified:      church.IsVerified.Bool,
			PostCount:       int64(church.PostsCount.Int32),
			About:           church.About.String,
			Website:         church.Website.String,
			HeaderImageUrl:  church.HeaderImageUrl.String,
			Location: &pb.Location{
				Country:    church.Country,
				State:      church.State,
				City:       church.City,
				Address:    church.Address,
				PostalCode: church.Postalcode,
				Lga:        church.Lga,
			},
		},
	}, nil

}
