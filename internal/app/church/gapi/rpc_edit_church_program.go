package gapi

import (
	"context"
	"database/sql"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ChurchServer) EditChurchProgram(ctx context.Context, req *pb.CreateChurchProgramRequest) (*pb.CreateChurchProgramResponse, error) {
	var startTime, endTime time.Time
	var err error

	if req.GetProgramStartTime() != nil {
		startTime, err = ptypes.Timestamp(req.ProgramStartTime)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert start time: %v", err)
		}
	}

	if req.GetProgramEndTime() != nil {
		endTime, err = ptypes.Timestamp(req.ProgramEndTime)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert end time: %v", err)
		}
	}

	// 	ProgramStartTime: sql.NullTime{Valid: true, Time: startTime},
	// 	ProgramEndTime:   sql.NullTime{Valid: true, Time: endTime},

	err = s.store.UpdateChurchProgram(ctx, sqlc.UpdateChurchProgramParams{
		ID:               req.GetId(),
		ProgramType:      sql.NullString{String: req.GetProgramType(), Valid: req.GetProgramType() != ""},
		ProgramName:      sql.NullString{String: req.GetProgramName(), Valid: req.GetProgramName() != ""},
		ProgramDesc:      sql.NullString{String: req.GetProgramDesc(), Valid: req.GetProgramDesc() != ""},
		ProgramFreq:      sql.NullString{String: req.GetProgramFreq(), Valid: req.GetProgramFreq() != ""},
		ProgramDay:       sql.NullString{String: req.GetProgramDay(), Valid: req.GetProgramDay() != ""},
		ProgramStartTime: sql.NullTime{Time: startTime, Valid: startTime != time.Time{}},
		ProgramEndTime:   sql.NullTime{Time: endTime, Valid: endTime != time.Time{}},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update program in the database: %v", err)
	}

	return &pb.CreateChurchProgramResponse{Msg: "Program updated successfully"}, nil

}
