package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"github.com/steve-mir/diivix_backend/constants"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	auth_service "github.com/steve-mir/diivix_backend/internal/app/authentication/services"
	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"github.com/steve-mir/diivix_backend/utils"
	"github.com/steve-mir/diivix_backend/worker"
	"golang.org/x/sync/errgroup"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidateCreateChurchInputs(req *pb.Church) error {
	if violations := validateCreateChurchRequest(req); violations != nil {
		return invalidArgumentErr(violations)
	}
	return nil
}

func CreateChurchAccount(ctx context.Context, req *pb.Church, db *sql.DB, store *sqlc.Store, config utils.Config, td worker.TaskDistributor) (*pb.CreateChurchResponse, error) {
	// get metadata
	agent, clientIP := utils.GetMetadata(ctx)

	// Begin transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := store.WithTx(tx)

	// Check db if user exists
	if err := auth_service.CheckUserExists(ctx, qtx, req.GetEmail(), req.GetUsername()); err != nil {
		return nil, status.Errorf(codes.AlreadyExists, err.Error())
	}

	// hash pwd and generate uuid
	hashedPwd, uid, err := auth_service.PrepareUserData(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error preparing data %+v", err)
	}

	err = RunConcurrentUserCreationTasks(ctx, qtx, tx, config, td, req, req.Location, uid, clientIP, hashedPwd, agent)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error running concurrent %+v", err)
	}

	return &pb.CreateChurchResponse{
		Success: true,
		Message: "Account created successfully. Please verify your email while you await verification by us",
	}, nil
}

func validateCreateChurchRequest(req *pb.Church) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := utils.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if !utils.ValidateUsername(req.GetUsername()) {
		violations = append(violations, fieldViolation("username", errors.New("invalid username format")))
	}

	if !utils.ValidatePhone(req.GetPhone()) {
		violations = append(violations, fieldViolation("phone", errors.New("invalid phone format")))
	}

	if !utils.ValidateAddress(req.Location.GetAddress()) {
		violations = append(violations, fieldViolation("address", errors.New("invalid address format")))
	}

	if !utils.ValidateCity(req.Location.GetCity()) {
		violations = append(violations, fieldViolation("city", errors.New("invalid city format")))
	}

	if !utils.ValidatePostalCode(req.Location.GetPostalCode()) {
		violations = append(violations, fieldViolation("postal code", errors.New("invalid postal code format")))
	}

	if !utils.ValidateState(req.Location.GetState()) {
		violations = append(violations, fieldViolation("state", errors.New("invalid state format")))
	}

	if !utils.ValidateCountry(req.GetLocation().GetCountry()) {
		violations = append(violations, fieldViolation("country", errors.New("invalid country format")))
	}

	if !utils.ValidatePassword(req.GetPassword()) {
		violations = append(violations, fieldViolation("password", errors.New("invalid password format")))
	}

	return violations
}

func RunConcurrentUserCreationTasks(ctx context.Context, qtx *sqlc.Queries, tx *sql.Tx, config utils.Config, td worker.TaskDistributor,
	req *pb.Church, location *pb.Location, uid uuid.UUID, clientIP pqtype.Inet, pwd, agent string) error {

	// Removed the wait groups for each goroutine and the individual channels
	// Instead, use an error group to manage concurrency
	var eg errgroup.Group

	// Variables to hold the results from goroutines
	// var auth sqlc.Authentication
	// var accessToken string
	// var accessPayload *token.Payload
	var err error

	// ! 1 Create Account
	_, err = auth_service.CreateNewUser(ctx, qtx, uid, req.GetEmail(), req.GetUsername(), pwd, constants.ChurchAdminUser)
	if err != nil {
		return err
	}

	// ! 2 Create access token (Should wait until the user is created)
	eg.Go(func() error {
		tokenService := auth_service.NewTokenService(config)
		_, _, err = tokenService.CreateAccessToken(req.GetEmail(), req.GetUsername(), req.GetPhone(), true, false, uid, constants.ChurchAdmin, clientIP, agent)
		if err != nil {
			return fmt.Errorf("error creating access token: %v", err)
		}
		return nil
	})

	// ! 3 Create user role (Can be concurrent)
	eg.Go(func() error {
		_, err = qtx.CreateUserRole(ctx, sqlc.CreateUserRoleParams{
			UserID: uid,
			RoleID: constants.ChurchAdmin,
		})
		if err != nil {
			return fmt.Errorf("error creating user role: %v", err)
		}
		return nil
	})

	// ! 4 Create church profile (Can be concurrent)
	eg.Go(func() error {
		err = qtx.CreateEntityProfile(ctx, sqlc.CreateEntityProfileParams{
			UserID:     uid,
			EntityType: constants.ChurchAdminUser,
		})
		if err != nil {
			return fmt.Errorf("error creating church profile: %v", err)
		}
		return nil
	})

	// ! 5 Create church data (Can be concurrent)
	eg.Go(func() error {
		err = qtx.CreateNewChurch(ctx, sqlc.CreateNewChurchParams{
			AuthID:         uid,
			DenominationID: req.GetDenominationId(),
			Name:           req.GetName(),
		})
		if err != nil {
			return fmt.Errorf("error creating church data: %v", err)
		}
		return nil
	})

	// ! 6 create church location (Can be concurrent)
	eg.Go(func() error {
		lat, lng, err := getLatAndLng(config, location.Address, location.City, location.PostalCode, location.State, location.Country)
		if err != nil {
			// Don't log here, return the error instead
			return fmt.Errorf("unable to get lat and lng: %v", err)
		}
		err = qtx.CreateNewChurchLocation(ctx, sqlc.CreateNewChurchLocationParams{
			AuthID:     uid,
			Address:    location.GetAddress(),
			City:       location.GetCity(),
			State:      location.GetState(),
			Postalcode: location.GetPostalCode(),
			Country:    location.GetCountry(),
			Lga:        location.GetLga(),
			Location:   fmt.Sprintf("SRID=4326;POINT(%f %f)", lat, lng), // Assuming SRID 4326 for WGS 84
		})
		if err != nil {
			return fmt.Errorf("error creating church location: %v", err)
		}
		return nil
	})

	// ! 7 Send verification email
	eg.Go(func() error {
		err := auth_service.SendVerificationEmail(qtx, ctx, td, uid, req.GetEmail())
		return err
	})

	// Wait for all goroutines to complete and return the first non-nil error (if any)
	if err := eg.Wait(); err != nil {
		tx.Rollback()
		return err
	}

	// Assuming accessPayload contains the expiry for the access token
	// expiryTime := time.Unix(accessPayload.Expires.Unix(), 0)

	// Return the results
	return tx.Commit()
}

// Makes network call to get coordinates
func getLatAndLng(config utils.Config, address, city, postalCode, state, country string) (float64, float64, error) {
	lat, lng, err := utils.TomTomGeocoding(config, address, city, postalCode, state, country)
	if err != nil {
		return 0, 0, fmt.Errorf("error getting location's coordinates %v", err)
	}
	return lat, lng, err
}
