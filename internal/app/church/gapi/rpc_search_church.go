package gapi

import (
	"context"
	"database/sql"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Default page size if not specified or negative
const DefaultPageSize = 10

func (s *ChurchServer) SearchChurches(ctx context.Context, req *pb.SearchRequestChurch) (*pb.SearchChurchResponse, error) {
	if req.GetQuery() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "query cannot be empty")
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

	churches, err := s.store.SearchChurches(ctx, sqlc.SearchChurchesParams{
		Column1: sql.NullString{String: req.GetQuery(), Valid: true},
		Offset:  int32(offset),
		Limit:   int32(pageSize),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error fetching churches: %s", err.Error())
	}

	resp := &pb.SearchChurchResponse{
		Church:  []*pb.Church{},
		HasMore: len(churches) == int(pageSize),
	}

	for _, church := range churches {
		// imageUrls := []string{}
		// if err := json.Unmarshal(church.PostImageUrls, &imageUrls); err != nil {
		// 	log.Printf("cannot unmarshal image URLs: %v", err)
		// 	continue
		// }

		// Skip posts that are deleted or suspended
		// if post.PostDeletedAt.Valid || post.PostSuspendedAt.Valid {
		// 	continue
		// }

		resp.Church = append(resp.Church, &pb.Church{
			Name:     church.Name,
			Username: church.Postalcode,
		})
	}

	return resp, nil
}
