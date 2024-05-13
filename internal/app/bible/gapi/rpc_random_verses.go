package gapi

import (
	"context"

	"github.com/steve-mir/diivix_backend/internal/app/bible/pb"
	"github.com/steve-mir/diivix_backend/internal/app/bible/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *BibleServer) RandomVerseGenerator(_ context.Context, req *pb.RandomVerseGeneratorRequest) (*pb.Verse, error) {
	book, chapter, verse, word := services.RandomVerseGenerator(req.GetBiblePath(), false)
	if word == "" {
		return nil, status.Errorf(codes.Internal, "could not fetch scripture")
	}
	return &pb.Verse{
		Book:    book,
		Chapter: int32(chapter),
		Verse:   int32(verse),
		Text:    word,
	}, nil
}
