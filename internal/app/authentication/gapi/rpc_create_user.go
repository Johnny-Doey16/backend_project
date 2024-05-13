package gapi

import (
	"context"
	"errors"
	"log"

	"github.com/steve-mir/diivix_backend/internal/app/authentication/pb"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/services"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if violations := validateCreateUserRequest(req); violations != nil {
		return nil, invalidArgumentErr(violations)
	}

	agent := server.extractMetadata(ctx).UserAgent
	ip := server.extractMetadata(ctx).ClientIP
	clientIP := utils.GetIpAddr(ip)

	// Begin transaction
	tx, err := server.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := server.store.WithTx(tx)

	// Check db if user exists
	if err := services.CheckUserExists(ctx, qtx, req.GetEmail(), req.GetUsername()); err != nil {
		return nil, status.Errorf(codes.AlreadyExists, err.Error())
	}

	// hash pwd and generate uuid
	hashedPwd, uid, err := services.PrepareUserData(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	sqlcUser, err := services.CreateUserConcurrent(ctx, qtx, tx, uid, req.Email, req.Username, hashedPwd)
	if err != nil {
		// tx.Rollback()
		return nil, status.Errorf(codes.Internal, "error while creating user with email and password %s", err)
	}

	// Run concurrent operations
	accessToken, accessExp, err := services.RunConcurrentUserCreationTasks(ctx, qtx, tx, server.config, server.taskDistributor, req, uid, clientIP, agent)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating details: %s", err.Error())
	}

	// Only commit the transaction if all previous steps were successful
	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "an unexpected error occurred during transaction commit: %s", err)
	}

	// TODO: Add username to cache: service.StoreUser()
	return &pb.CreateUserResponse{
		User: &pb.User{
			Uid:             uid.String(),
			IsEmailVerified: sqlcUser.IsEmailVerified.Bool,
			IsVerified:      sqlcUser.IsVerified.Bool,
			Username:        sqlcUser.Username.String,
			Email:           sqlcUser.Email,
			CreatedAt:       timestamppb.New(sqlcUser.CreatedAt.Time),
		},
		AccessToken:          accessToken,
		AccessTokenExpiresAt: timestamppb.New(accessExp),
	}, nil

}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := utils.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if !utils.ValidateUsername(req.GetUsername()) {
		violations = append(violations, fieldViolation("username", errors.New("invalid username format")))
	}

	if !utils.ValidatePassword(req.GetPassword()) {
		log.Println("PWD", req.GetPassword())
		violations = append(violations, fieldViolation("password", errors.New("invalid password format")))
	}

	return violations
}
