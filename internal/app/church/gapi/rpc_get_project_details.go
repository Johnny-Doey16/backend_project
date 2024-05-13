package gapi

import (
	"context"

	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ChurchServer) GetProjectDetails(ctx context.Context, req *pb.GetProjectDetailsRequest) (*pb.ChurchProject, error) {
	project, err := s.store.GetChurchProjectDetails(ctx, int32(req.GetProjectId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "SQL ERROR: "+err.Error())
	}
	// fmt.Println("User Details type: ", reflect.TypeOf(project[0].Contributors))
	// var contrib []*pb.ProjectContributor
	// err = json.Unmarshal(project[0].Contributors, &contrib)
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "failed to unmarshal user details: %v", err)
	// }

	return &pb.ChurchProject{
		ChurchId:    int32(project[0].ChurchID),
		Id:          int32(project[0].ID),
		Name:        project[0].ProjectName,
		Description: project[0].ProjectDescription.String,
		Visibility:  project[0].Visibility.Bool,
		Completed:   project[0].Completed.Bool,
		StartDate:   timestamppb.New(project[0].StartDate.Time),
		EndDate:     timestamppb.New(project[0].EndDate.Time),
		// TargetAmount: project[0].TargetAmount,
		// DonatedAmount: project[0].DonatedAmount,

		// Contributors: contrib,
	}, nil
}
