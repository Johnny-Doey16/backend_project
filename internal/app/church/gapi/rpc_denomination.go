package gapi

import (
	"context"

	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ChurchServer) AddDenomination(ctx context.Context, req *pb.AddDenominationRequest) (*pb.AddDenominationResponse, error) {
	if req.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Argument cannot be empty")
	}

	if len(req.GetName()) < 5 {
		return nil, status.Errorf(codes.InvalidArgument, "error adding denomination")
	}

	if err := s.store.CreateDenomination(ctx, req.GetName()); err != nil {
		return nil, status.Errorf(codes.Internal, "error creating denominations %v", err.Error())
	}

	return &pb.AddDenominationResponse{
		Message: "Successfully created new denomination",
		Success: true,
	}, nil

}

func (s *ChurchServer) GetDenominationList(ctx context.Context, req *pb.GetDenominationRequest) (*pb.GetDenominationResponse, error) {
	denominations, err := s.store.GetDenominationList(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error fetching denominations %v", err.Error())
	}

	resp := &pb.GetDenominationResponse{
		Denominations: []*pb.Denomination{},
	}

	for _, denomination := range denominations {
		resp.Denominations = append(resp.Denominations, &pb.Denomination{
			Id:   denomination.ID,
			Name: denomination.Name,
		})
	}

	return resp, nil
}
