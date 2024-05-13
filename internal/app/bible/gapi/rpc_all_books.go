package gapi

import (
	"github.com/steve-mir/diivix_backend/internal/app/bible/pb"
	"github.com/steve-mir/diivix_backend/internal/app/bible/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *BibleServer) ListAllBooks(req *pb.ListAllBooksRequest, stream pb.BibleService_ListAllBooksServer) error {
	books, err := services.ListAllBooksOfBible(req.GetBiblePath()) // "bible-translations/eng-dra.osis.xml"
	if err != nil {
		return err
	}

	// Stream each chapter to the client
	for _, book := range books {
		resBook := &pb.ListAllBooksResponse{
			BookName: book,
		}

		if err := stream.Send(resBook); err != nil {
			return status.Errorf(codes.Internal, "error sending all books of the bible. Error: %s", err)
		}
	}

	return nil
}
