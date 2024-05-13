package gapi

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/constants"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

/*
func (s *SocialMediaServer) ViewUserProfile(ctx context.Context, req *pb.UserProfileRequest) (*pb.UserProfileResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}
	var user sqlc.GetUserProfileRow
	var err error

	switch identifier := req.GetIdentifier().(type) {
	//! Is UID
	case *pb.UserProfileRequest_Uid:
		userProfileId, err := services.StrToUUID(identifier.Uid)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not convert string to uid %v", err)
		}

		user, err = s.store.GetUserProfile(ctx, sqlc.GetUserProfileParams{
			FollowerUserID: claims.UserId, // The current user id
			ID:             userProfileId, // visiting user id
		})

		//! Is Username
	case *pb.UserProfileRequest_Username:
		user2, err = s.store.GetUserProfileByUsername(ctx, sqlc.GetUserProfileByUsernameParams{
			FollowerUserID: claims.UserId,
			Username:       sql.NullString{Valid: true, String: req.GetUid()},
		})

	}





	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.Internal, "User blocked you. Cannot view user profile")
		}
		return nil, status.Errorf(codes.Internal, "could not find the user %v", err)
	}

	return &pb.UserProfileResponse{
		User: &pb.PostUser{
			Uid:            user.ID.String(),
			Username:       user.Username.String,
			Email:          user.Email,
			Phone:          user.Phone.String,
			FirstName:      user.FirstName.String,
			LastName:       user.LastName.String,
			ImageUrl:       user.ImageUrl.String,
			IsVerified:     user.IsVerified.Bool,
			FollowingCount: user.FollowingCount.Int32,
			FollowerCount:  user.FollowersCount.Int32,
			CreatedAt:      timestamppb.New(user.CreatedAt.Time),
		},
		IsFollowing: user.IsFollowing,
		IsFollowed:  user.IsFollowed,
	}, nil
}
*/

func (s *SocialMediaServer) ViewUserProfile(ctx context.Context, req *pb.UserProfileRequest) (*pb.UserProfileResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	/*

		var err error
		var userProfile *pb.UserProfileResponse

		switch identifier := req.GetIdentifier().(type) {
		case *pb.UserProfileRequest_Uid:
			fmt.Println("Getting by UID")
			userProfile, err = s.getUserProfileByUID(ctx, claims.UserId, identifier.Uid)

		case *pb.UserProfileRequest_Username:
			fmt.Println("Getting by username")
			userProfile, err = s.getUserProfileByUsername(ctx, claims.UserId, identifier.Username)

		default:
			return nil, status.Error(codes.InvalidArgument, "invalid identifier type")
		}
	*/

	userProfile, err := s.getUserProfileByUID(ctx, claims.UserId, req.GetIdentifier())
	if err != nil {
		return nil, handleUserProfileError(err)
	}

	return userProfile, nil
}

func (s *SocialMediaServer) getUserProfileByUID(ctx context.Context, followerUserID uuid.UUID, uid string) (*pb.UserProfileResponse, error) {
	// userProfileId, err := services.StrToUUID(uid)
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "could not convert string to UID: %v", err)
	// }

	user, err := s.store.GetUserProfile(ctx, sqlc.GetUserProfileParams{ //GetUserProfile
		FollowerUserID: followerUserID,
		Username:       sql.NullString{Valid: true, String: uid},
	})

	if err != nil {
		return nil, err
	}

	return buildUserProfileResponse(user), nil
}

func handleUserProfileError(err error) error {
	if err == sql.ErrNoRows {
		return status.Errorf(codes.Internal, "User blocked you. Cannot view user profile: %v", err)
	}
	return status.Errorf(codes.Internal, "could not find the user: %v", err)
}

func buildUserProfileResponse(user sqlc.GetUserProfileRow) *pb.UserProfileResponse {
	var firstName *string
	var lastName *string
	var church *pb.ChurchUser
	var _user *pb.PostUser
	var isChurch bool

	if user.UserType == constants.ChurchAdminUser {
		isChurch = true
		church = &pb.ChurchUser{
			Uid:             user.ID.String(),
			Id:              user.ChurchID.Int32,
			Email:           user.Email,
			Username:        user.Username.String,
			Phone:           user.Phone.String,
			DenominationId:  user.ChurchDenominationID.Int32,
			Name:            *services.InterfaceToStr(user.FirstName, firstName),
			About:           user.About.String,
			Website:         user.Website.String,
			HeaderImageUrl:  user.HeaderImageUrl.String,
			PostCount:       int64(user.PostsCount.Int32),
			CreatedAt:       timestamppb.New(user.CreatedAt.Time),
			ImageUrl:        user.ImageUrl.String,
			IsVerified:      user.IsVerified.Bool,
			FollowingCount:  int64(user.FollowingCount.Int32),
			FollowerCount:   int64(user.FollowersCount.Int32),
			MembershipCount: int64(user.ChurchMembersCount.Int32),
			AccountName:     user.AccountName.String,
			AccountNumber:   user.AccountNumber.String,
			BankName:        user.BankName.String,
			Location: &pb.LocationUser{
				Address:    user.ChurchAddress.String,
				City:       user.ChurchCity.String,
				PostalCode: user.ChurchPostalcode.String,
				State:      user.ChurchState.String,
				Country:    user.ChurchCountry.String,
				Lga:        user.ChurchLga.String,
			},
		}
	}

	if user.UserType == constants.RegularUsersUser {
		isChurch = false
		_user = &pb.PostUser{
			Uid:            user.ID.String(),
			Username:       user.Username.String,
			Email:          user.Email,
			Phone:          user.Phone.String,
			FirstName:      *services.InterfaceToStr(user.FirstName, firstName),
			LastName:       *services.InterfaceToStr(user.LastName, lastName),
			ImageUrl:       user.ImageUrl.String,
			IsVerified:     user.IsVerified.Bool,
			FollowingCount: user.FollowingCount.Int32,
			FollowerCount:  user.FollowersCount.Int32,
			CreatedAt:      timestamppb.New(user.CreatedAt.Time),
			About:          user.About.String,
			Website:        user.Website.String,
			HeaderImageUrl: user.HeaderImageUrl.String,
			PostCount:      user.PostsCount.Int32,
		}
	}

	return &pb.UserProfileResponse{
		User:        _user,
		Church:      church,
		IsFollowing: user.IsFollowing,
		IsFollowed:  user.IsFollowed,
		IsMember:    user.IsMember,
		IsChurch:    isChurch,
	}
}

/*
func (s *SocialMediaServer) getUserProfileByUsername(ctx context.Context, followerUserID uuid.UUID, username string) (*pb.UserProfileResponse, error) {
	user, err := s.store.GetUserProfileByUsername(ctx, sqlc.GetUserProfileByUsernameParams{
		FollowerUserID: followerUserID,
		Username:       sql.NullString{Valid: true, String: username},
	})

	if err != nil {
		return nil, err
	}

	return buildUserProfileResponseByUsername(user), nil
}
func buildUserProfileResponseByUsername(user sqlc.GetUserProfileByUsernameRow) *pb.UserProfileResponse {
	// Adjust the logic based on the actual structure of GetUserProfileByUsernameRow
	// Assuming it's similar to GetUserProfileRow
	var firstName *string
	var lastName *string
	var church *pb.ChurchUser
	var _user *pb.PostUser

	if user.UserType == constants.ChurchAdminUser {
		church = &pb.ChurchUser{
			AuthId:         user.ID.String(),
			Id:             user.ChurchID.Int32,
			Email:          user.Email,
			Username:       user.Username.String,
			Phone:          user.Phone.String,
			DenominationId: user.ChurchDenominationID.Int32,
			Name:           *services.InterfaceToStr(user.FirstName, firstName),
			About:          user.About.String,
			Website:        user.Website.String,
			HeaderImageUrl: user.HeaderImageUrl.String,
			PostCount:      int64(user.PostsCount.Int32),
			CreatedAt:      timestamppb.New(user.CreatedAt.Time),
			ImageUrl:       user.ImageUrl.String,
			IsVerified:     user.IsVerified.Bool,
			FollowingCount: int64(user.FollowingCount.Int32),

			FollowerCount:   int64(user.FollowersCount.Int32),
			MembershipCount: int64(user.ChurchMembersCount.Int32),
			Location: &pb.LocationUser{
				Address:    user.ChurchAddress.String,
				City:       user.ChurchCity.String,
				PostalCode: user.ChurchPostalcode.String,
				State:      user.ChurchState.String,
				Country:    user.ChurchCountry.String,
				Lga:        user.ChurchLga.String,
			},
		}
	}

	if user.UserType == constants.RegularUsersUser {
		_user = &pb.PostUser{
			Uid:            user.ID.String(),
			Username:       user.Username.String,
			Email:          user.Email,
			Phone:          user.Phone.String,
			FirstName:      *services.InterfaceToStr(user.FirstName, firstName),
			LastName:       *services.InterfaceToStr(user.LastName, lastName),
			ImageUrl:       user.ImageUrl.String,
			IsVerified:     user.IsVerified.Bool,
			FollowingCount: user.FollowingCount.Int32,
			FollowerCount:  user.FollowersCount.Int32,
			CreatedAt:      timestamppb.New(user.CreatedAt.Time),
			About:          user.About.String,
			Website:        user.Website.String,
			HeadImageUrl:   user.HeaderImageUrl.String,
			PostCount:      user.PostsCount.Int32,
		}
	}

	return &pb.UserProfileResponse{
		User:        _user,
		Church:      church,
		IsFollowing: user.IsFollowing,
		IsFollowed:  user.IsFollowed,
		IsMember:    user.IsMember,
	}
}

*/
