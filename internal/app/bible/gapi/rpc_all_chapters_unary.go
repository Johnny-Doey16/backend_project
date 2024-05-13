package gapi

import (
	"context"

	"github.com/steve-mir/diivix_backend/internal/app/bible/pb"
	"github.com/steve-mir/diivix_backend/internal/app/bible/services"
)

func (s *BibleServer) ListAllChaptersOfBookUnary(ctx context.Context, req *pb.ListAllChaptersOfBookRequest) (*pb.ListAllChaptersOfBookUnaryResponse, error) {
	chapters, err := services.ListAllChaptersOfBook(req.GetBiblePath(), req.GetBookName())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListAllChaptersOfBookUnaryResponse{
		Chapter: []int32{},
	}

	for _, chapter := range chapters {
		resp.Chapter = append(resp.Chapter, int32(chapter))
	}

	return resp, nil
}
