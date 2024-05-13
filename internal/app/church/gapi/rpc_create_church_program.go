package gapi

import (
	"context"
	"database/sql"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"github.com/steve-mir/diivix_backend/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ChurchServer) CreateProgram(ctx context.Context, req *pb.CreateChurchProgramRequest) (*pb.CreateChurchProgramResponse, error) {
	sTime := utils.StartMemCal()

	// TODO: Add church ID to
	// Convert Timestamps to Go Time
	programStartTime, err := ptypes.Timestamp(req.GetProgramStartTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert start time: %v", err)
	}
	programEndTime, err := ptypes.Timestamp(req.GetProgramEndTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert end time: %v", err)
	}

	err = s.store.CreateChurchProgram(ctx, sqlc.CreateChurchProgramParams{
		ChurchID:         15,
		ProgramType:      req.GetProgramType(),
		ProgramName:      req.GetProgramName(),
		ProgramDesc:      sql.NullString{String: req.GetProgramDesc(), Valid: req.GetProgramDesc() != ""},
		ProgramFreq:      sql.NullString{String: req.GetProgramFreq(), Valid: req.GetProgramFreq() != ""},
		ProgramDay:       sql.NullString{String: req.GetProgramDay(), Valid: req.GetProgramDay() != ""},
		ProgramStartTime: sql.NullTime{Valid: programStartTime != time.Time{}, Time: programStartTime},
		ProgramEndTime:   sql.NullTime{Valid: programEndTime != time.Time{}, Time: programEndTime},
		// ProgramImageUrl: sql.NullString{String: "", Valid: false},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert program into the database: %v", err)
	}

	utils.EndMemCal(sTime)

	return &pb.CreateChurchProgramResponse{Msg: "Program created successfully"}, nil

}

// Peak Memory Usage: 0.0125 GB
// Execution Time: 30.1127 seconds
// GB-Seconds: 0.3750
