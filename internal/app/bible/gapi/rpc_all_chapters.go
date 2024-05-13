package gapi

import (
	"github.com/steve-mir/diivix_backend/internal/app/bible/pb"
	"github.com/steve-mir/diivix_backend/internal/app/bible/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *BibleServer) ListAllChaptersOfBook(req *pb.ListAllChaptersOfBookRequest, stream pb.BibleService_ListAllChaptersOfBookServer) error {
	chapters, err := services.ListAllChaptersOfBook(req.GetBiblePath(), req.GetBookName())
	if err != nil {
		return err
	}

	// Stream each chapter to the client
	for _, chapter := range chapters {
		resChapter := &pb.ListAllChaptersOfBookResponse{
			Chapter: int32(chapter),
		}

		if err := stream.Send(resChapter); err != nil {
			return status.Errorf(codes.Internal, "error sending chapters of %s. Error: %s", req.GetBookName(), err)
		}
	}

	return nil
}
