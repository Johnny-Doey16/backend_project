package gapi

import (
	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ChurchServer) GetChurchProjects(req *pb.GetChurchProjectsRequest, stream pb.ChurchService_GetChurchProjectsServer) error {
	projects, err := s.store.GetChurchProjects(stream.Context(), int32(req.GetChurchId()))
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for _, project := range projects {
		stream.Send(&pb.GetChurchProjectsResponse{
			Id:          int32(project.ID),
			ChurchId:    int32(project.ChurchID),
			Name:        project.ProjectName,
			Description: project.ProjectDescription.String,
			Visibility:  project.Visibility.Bool,
			Completed:   project.Completed.Bool,
			StartDate:   timestamppb.New(project.StartDate.Time),
			EndDate:     timestamppb.New(project.EndDate.Time),
			// TargetAmount: float64(project.TargetAmount),
			// DonatedAmount: float64(project.DonatedAmount),

		})
	}
	/*
	   bool visibility = 9;
	   bool completed = 10;
	*/
	return nil
}
