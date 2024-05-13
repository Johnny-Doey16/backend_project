package gapi

import (
	"context"
	"database/sql"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO: Fix data input type issue
func (s *ChurchServer) CreateProject(ctx context.Context, req *pb.CreateProjectRequest) (*pb.CreateProjectResponse, error) {
	err := s.store.CreateChurchProject(ctx, sqlc.CreateChurchProjectParams{
		ChurchID:           int32(req.GetChurchId()),
		ProjectName:        req.GetProjectName(),
		ProjectDescription: sql.NullString{String: req.GetProjectDescription(), Valid: true},
		// TargetAmount:       req.GetTargetAmount(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.CreateProjectResponse{
		Msg: "Successfully created",
	}, nil
}
