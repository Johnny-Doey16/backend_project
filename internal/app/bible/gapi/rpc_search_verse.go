package gapi

import (
	"context"

	"github.com/steve-mir/diivix_backend/internal/app/bible/pb"
	"github.com/steve-mir/diivix_backend/internal/app/bible/services"
)

func (s *BibleServer) SearchVerse(ctx context.Context, req *pb.SearchVerseRequest) (*pb.Verse, error) {
	book, chapter, verse, text, err := services.SearchBibleVerse(req.GetBiblePath(), req.GetInput())
	if err != nil {
		return nil, err
	}
	return &pb.Verse{
		Book:    book,
		Chapter: int32(chapter),
		Verse:   int32(verse),
		Text:    text,
	}, nil
}
