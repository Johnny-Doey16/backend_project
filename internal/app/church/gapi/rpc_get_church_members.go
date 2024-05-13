package gapi

import (
	"encoding/json"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetChurchMembers(ctx context.Context, req *pb.GetChurchMembersRequest) (*pb.GetChurchMembersResponse, error) {
func (s *ChurchServer) GetChurchMembers(req *pb.GetChurchMembersRequest, stream pb.ChurchService_GetChurchMembersServer) error {
	claims, ok := stream.Context().Value("payloadKey").(*token.Payload)
	if !ok {
		return status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	pageNumber := req.GetPageNumber()
	if pageNumber < 1 {
		pageNumber = 1
	}

	pageSize := req.GetPageSize()
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}

	offset := (pageNumber - 1) * pageSize

	// if strings.ToLower(req.GetOrder()) != "ASC" || strings.ToLower(req.GetOrder()) != "DESC" {
	// 	return nil, status.Error(codes.InvalidArgument, "order should be either 'asc' or 'desc'")
	// }

	members, err := s.store.GetChurchMembers5(stream.Context(), sqlc.GetChurchMembers5Params{
		ChurchID: req.GetChurchId(),
		UserID:   claims.UserId,
		Offset:   offset,
		Limit:    pageSize,
		// Column3:  sql.NullString{String: req.GetOrder(), Valid: true},
	})
	if err != nil {
		return status.Error(codes.Internal, "cannot fetch members data "+err.Error())
	}

	// var churchMembers []*pb.ChurchMember

	var userDetails []pb.ChurchMember
	err = json.Unmarshal(members[0].UserDetails, &userDetails)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to unmarshal user details: %v", err)
	}

	for i := range userDetails {
		// Use a pointer to the original struct
		// churchMember := &pb.ChurchMember{
		// 	AuthId:     userDetails[i].AuthId,
		// 	ImageUrl:   userDetails[i].ImageUrl,
		// 	FirstName:  userDetails[i].FirstName,
		// 	Username:   userDetails[i].Username,
		// 	IsVerified: userDetails[i].IsVerified,
		// }
		// churchMembers = append(churchMembers, churchMember)

		stream.Send(&pb.GetChurchMembersResponse{
			IsMember: members[0].IsMember,
			Members: &pb.ChurchMember{
				AuthId:     userDetails[i].AuthId,
				ImageUrl:   userDetails[i].ImageUrl,
				FirstName:  userDetails[i].FirstName,
				Username:   userDetails[i].Username,
				IsVerified: userDetails[i].IsVerified,
			},
		})
	}

	// return &pb.GetChurchMembersResponse{
	// 	Members:  churchMembers,
	// 	IsMember: members[0].IsMember,
	// }, nil
	return nil

}
