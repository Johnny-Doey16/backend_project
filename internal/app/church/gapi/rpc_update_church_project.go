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
func (s *ChurchServer) UpdateProject(ctx context.Context, req *pb.UpdateProjectRequest) (*pb.UpdateProjectResponse, error) {
	err := s.store.UpdateChurchProject(ctx, sqlc.UpdateChurchProjectParams{
		ID:                 int32(req.GetProjectId()),
		ProjectName:        sql.NullString{String: req.GetName(), Valid: req.Name != ""},
		ProjectDescription: sql.NullString{String: req.GetDescription(), Valid: req.Description != ""},
		Visibility:         sql.NullBool{Bool: req.GetVisibility(), Valid: req.Visibility != false},
		// TargetAmount: sql.NullString{String: req.GetTargetAmount(), Valid: req.TargetAmount != 0},
		// ProjectName: sql.NullTime{String: req.GetName(), Valid: req.Name != ""},
		// EndDate: sql.NullTime{String: req.GetEndDate(), Valid: req.EndDate != ""},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.UpdateProjectResponse{
		Msg: "Successfully created",
	}, nil
}
