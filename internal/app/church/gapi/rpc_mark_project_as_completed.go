package gapi

import (
	"context"

	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ChurchServer) MarkProjectCompleted(ctx context.Context, req *pb.MarkProjectCompletedRequest) (*pb.MarkProjectCompletedResponse, error) {
	err := s.store.MarkChurchProjectAsCompleted(ctx,
		int32(req.GetProjectId()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.MarkProjectCompletedResponse{
		Msg: "Successfully created",
	}, nil
}
