package gapi

import (
	"log"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"github.com/steve-mir/diivix_backend/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ChurchServer) GetChurchPrograms(req *pb.GetChurchProgramsRequest, stream pb.ChurchService_GetChurchProgramsServer) error {
	sTime := utils.StartMemCal()
	log.Println("Calling with request: ", req)

	pageNumber := req.PageNumber
	pageSize := req.PageSize
	if pageNumber < 1 {
		pageNumber = 1
	}
	if pageSize < 1 {
		pageSize = 10 // Default page size to 10 if not specified or negative
	}

	// Calculate the offset based on the page number and page size.
	offset := (pageNumber - 1) * pageSize

	programs, err := s.store.GetChurchProgramsByChurchId(stream.Context(), sqlc.GetChurchProgramsByChurchIdParams{
		ChurchID: int32(req.GetChurchId()),
		Limit:    pageSize,
		Offset:   int32(offset),
	})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to fetch church programs: %v", err)
	}
	log.Println("Length of programs: ", len(programs))

	for _, program := range programs {
		stream.Send(&pb.GetChurchProgramsResponse{
			ChurchProgram: &pb.CreateChurchProgramRequest{
				Id:               program.ID,
				ProgramType:      program.ProgramType,
				ProgramName:      program.ProgramName,
				ProgramDesc:      program.ProgramDesc.String,
				ProgramDay:       program.ProgramDay.String,
				ProgramStartTime: timestamppb.New(program.ProgramStartTime.Time),
				ProgramEndTime:   timestamppb.New(program.ProgramEndTime.Time),
				ProgramFreq:      program.ProgramFreq.String,
			},
		})
	}

	utils.EndMemCal(sTime)
	return nil
}

// func (s *ChurchServer) GetChurchPrograms(req *pb.GetChurchProgramsRequest, stream pb.ChurchService_GetChurchProgramsServer) error {
// 	sTime := utils.StartMemCal()
// 	rows, err := s.db.Query(`
// 	    SELECT id, program_type, program_name, program_desc, program_day, program_start_time, program_end_time, program_freq
// 	    FROM church_programs
// 	    WHERE church_id = $1
// 	    OFFSET $2
// 	    LIMIT $3`,
// 		req.ChurchId, req.PageNumber*req.PageSize, req.PageSize)
// 	if err != nil {
// 		return status.Errorf(codes.Internal, "failed to fetch church programs: %v", err)
// 	}
// 	defer rows.Close()
// 	// programs, err := s.store.GetChurchProgramsByChurchId(req.GetChurchId())
// 	// if err != nil {
// 	// 	return status.Errorf(codes.Internal, "failed to fetch church programs: %v", err)
// 	// }

// 	for rows.Next() {
// 		var programType, programName, programDesc, programDay, programFreq string
// 		var id int32
// 		var startTime, endTime time.Time
// 		if err := rows.Scan(&id, &programType, &programName, &programDesc, &programDay, &startTime, &endTime, &programFreq); err != nil {
// 			return status.Errorf(codes.Internal, "failed to scan row: %v", err)
// 		}

// 		// Convert time.Time to Timestamp
// 		startTimeProto := timestamppb.New(startTime)
// 		endTimeProto := timestamppb.New(endTime)

// 		// Send the program details through the stream
// 		if err := stream.Send(&pb.GetChurchProgramsResponse{
// 			ChurchProgram: &pb.CreateChurchProgramRequest{
// 				Id:               id,
// 				ProgramType:      programType,
// 				ProgramName:      programName,
// 				ProgramDesc:      programDesc,
// 				ProgramDay:       programDay,
// 				ProgramStartTime: startTimeProto,
// 				ProgramEndTime:   endTimeProto,
// 				ProgramFreq:      programFreq,
// 			},
// 		}); err != nil {
// 			return status.Errorf(codes.Internal, "failed to send program details: %v", err)
// 		}
// 	}

// 	if err := rows.Err(); err != nil {
// 		return status.Errorf(codes.Internal, "error encountered while iterating over rows: %v", err)
// 	}
// 	utils.EndMemCal(sTime)
// 	return nil
// }
