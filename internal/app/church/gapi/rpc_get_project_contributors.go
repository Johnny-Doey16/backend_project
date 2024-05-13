package gapi

import (
	"encoding/json"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ChurchServer) GetProjectContributors(req *pb.GetProjectDetailsRequest, stream pb.ChurchService_GetProjectContributorsServer) error {
	contribData, err := s.store.GetChurchProjectContributors(stream.Context(), sqlc.GetChurchProjectContributorsParams{
		ProjectID: int32(req.GetProjectId()),
		// Limit:     1,
		// Offset:    0,
	})
	if err != nil {
		return status.Error(codes.Internal, "SQL ERROR: "+err.Error())
	}

	var contrib []*pb.ProjectContributor
	err = json.Unmarshal(contribData[0], &contrib)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to unmarshal user details: %v", err)
	}

	for i := range contrib {
		// Use a pointer to the original struct
		stream.Send(&pb.ProjectContributor{
			UserId:        contrib[i].UserId,
			TotalDonation: contrib[i].TotalDonation,
		})
	}

	return nil
}
