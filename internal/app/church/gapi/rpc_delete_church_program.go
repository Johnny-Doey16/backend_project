package gapi

import (
	"context"

	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ChurchServer) DeleteChurchProgram(ctx context.Context, req *pb.DeleteChurchProgramRequest) (*pb.CreateChurchProgramResponse, error) {
	err := s.store.DeleteChurchProgram(ctx, req.ProgramId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete program from the database: %v", err)
	}
	return &pb.CreateChurchProgramResponse{Msg: "Program deleted successfully"}, nil
}
