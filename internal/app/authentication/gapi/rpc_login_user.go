package gapi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/pb"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/services"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	if violations := validateLoginUserRequest(req.GetIdentifier(), req.GetPassword()); violations != nil {
		return nil, invalidArgumentErr(violations)
	}

	agent := server.extractMetadata(ctx).UserAgent
	ip := server.extractMetadata(ctx).ClientIP

	clientIP := utils.GetIpAddr(ip)
	log.Println("User ip", clientIP, " Agent", agent)

	sessionID, err := uuid.NewRandom()
	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, "error creating uid %s", err)
	}

	user, err := server.store.GetUserAndRoleByIdentifier(ctx, sql.NullString{String: req.GetIdentifier(), Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.Unimplemented, "email or password incorrect")
		}
		return nil, status.Error(codes.Unimplemented, "email or password incorrect")
	}

	err = utils.CheckPassword(req.GetPassword(), user.PasswordHash)
	if err != nil {
		return nil, status.Error(codes.Unimplemented, "email or password incorrect")
	}

	// Check if user should gain access
	err = checkAccountStat(user.IsSuspended.Bool, user.IsDeleted.Bool)
	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, "checking account stat error %s", err)
	}

	// Check if MFA is enabled for the user
	// if user.IsMfaEnabled.Bool {
	// isMfaEnabled := true
	// return &pb.LoginUserResponse{
	// 	User: &pb.User{
	// 		IsMfaEnable: &isMfaEnabled,
	// 	},
	// }, nil
	// }

	var mfaPassed bool
	if user.IsMfaEnabled.Bool {
		mfaPassed = false
	} else {
		mfaPassed = true
	}

	tokenService := services.NewTokenService(server.config)
	// Refresh token
	refreshToken, refreshPayload, err := tokenService.CreateRefreshToken(user.ID, sessionID, clientIP, agent)
	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, "error creating refresh token %s", err)
	}

	// Access token
	accessToken, accessPayload, err := tokenService.CreateAccessToken(user.Email, user.Username.String, user.Phone.String, mfaPassed,
		user.IsEmailVerified.Bool, user.ID, int8(user.RoleID), clientIP, agent)
	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, "error creating access token %s", err)
	}

	_, err = server.store.CreateSession(ctx, sqlc.CreateSessionParams{
		ID:              sessionID,
		UserID:          user.ID,
		RefreshToken:    refreshToken,
		RefreshTokenExp: refreshPayload.Expires,
		UserAgent:       agent,
		IpAddress:       clientIP,
		FcmToken:        sql.NullString{String: req.GetFcmToken(), Valid: true},
	})

	if err != nil {
		log.Println("Session ID Error", err)
		return nil, status.Errorf(codes.Unimplemented, "error creating session id %s", err)
	}

	// ctx.SetCookie("accessToken", accessToken, 36000, "/", "http://localhost:9100/", false, true)

	//! 3 User logged in successfully. Record it
	err = recordLoginSuccess(ctx, server.store, user.ID, agent, clientIP)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating login record %s", err)
	}

	// return resp
	return &pb.LoginUserResponse{
		User: &pb.User{
			Uid:               user.ID.String(),
			IsEmailVerified:   user.IsEmailVerified.Bool,
			IsVerified:        user.IsVerified.Bool,
			IsDeleted:         user.IsDeleted.Bool,
			IsSuspended:       user.IsSuspended.Bool,
			Username:          user.Username.String,
			Email:             user.Email,
			CreatedAt:         timestamppb.New(user.CreatedAt.Time),
			PasswordChangedAt: timestamppb.New(user.PasswordLastChanged.Time),
			IsMfaEnable:       user.IsMfaEnabled.Bool,
			About:             user.About.String,
			ImageUrl:          user.ImageUrl.String,
			HeaderImageUrl:    user.HeaderImageUrl.String,
			Website:           user.Website.String,
			FollowerCount:     user.FollowersCount.Int32,
			FollowingCount:    user.FollowingCount.Int32,
			TotalCoinsString:  user.TotalCoin.String,
		},
		AuthInfo: &pb.AuthTokenInfo{
			AccessToken:           accessToken,
			AccessTokenExpiresAt:  timestamppb.New(accessPayload.Expires),
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: timestamppb.New(refreshPayload.Expires),
		},
	}, nil
}

// *********
func validateLoginUserRequest(identifier string, pwd string) (violations []*errdetails.BadRequest_FieldViolation) {
	if utils.IsEmailFormat(identifier) { // Assuming there's a function to check if the format is an email
		if err := utils.ValidateEmail(identifier); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	} else if utils.IsPhoneFormat(identifier) {
		if !utils.ValidatePhone(identifier) {
			violations = append(violations, fieldViolation("phone", errors.New("invalid phone format")))
		}
	} else { // Default to username validation if it's not email or phone
		if !utils.ValidateUsername(identifier) {
			violations = append(violations, fieldViolation("username", errors.New("invalid username format")))
		}
	}

	// Continue with password validation
	if !utils.ValidatePassword(pwd) {
		violations = append(violations, fieldViolation("password", errors.New("invalid password format")))
	}

	return violations
}

func checkAccountStat(isSuspended bool, isDeleted bool) error {
	fmt.Printf("Is Suspended %v is deleted %v", isSuspended, isDeleted)
	if isSuspended {
		log.Println("Account deleted: ", isSuspended)
		return errors.New("account suspended")
	}

	// Check if user should gain access
	if isDeleted {
		log.Println("Account deleted: ", isDeleted)
		return errors.New("account suspended")
	}
	return nil
}

func recordLoginSuccess(ctx context.Context, dbStore *sqlc.Store, userId uuid.UUID, userAgent string, ipAddrs pqtype.Inet) error {
	_, err := dbStore.CreateUserLogin(ctx, sqlc.CreateUserLoginParams{
		UserID: userId,
		UserAgent: sql.NullString{
			String: userAgent,
			Valid:  true,
		},
		IpAddress: ipAddrs,
	})
	return err
}
