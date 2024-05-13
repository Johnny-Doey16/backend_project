package gapi

import (
	"context"

	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"github.com/steve-mir/diivix_backend/internal/app/church/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ChurchServer) CreateChurch(ctx context.Context, req *pb.Church) (*pb.CreateChurchResponse, error) {

	/***
		string address = 1;
	    string city = 2;
	    string postalCode = 3;
	    string state = 4;
	    string country = 5;

		string email = 1;
		string username = 2;
		string phone = 3;
		string password = 4;
		int32 denomination_id = 5;
		string name = 6;
		Location location = 7;
	*/

	// Sanitize inputs
	if err := services.ValidateCreateChurchInputs(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	res, err := services.CreateChurchAccount(ctx, req, s.db, s.store, s.config, s.taskDistributor)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return res, nil
}
