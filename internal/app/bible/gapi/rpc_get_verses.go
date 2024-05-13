package gapi

import (
	"context"

	"github.com/steve-mir/diivix_backend/internal/app/bible/pb"
	"github.com/steve-mir/diivix_backend/internal/app/bible/services"
)

func (s *BibleServer) GetVerseRange(_ context.Context, req *pb.GetVerseRangeRequest) (*pb.VerseRange, error) {
	book, chapter, startVerse, endVerse, verseTexts, err := services.GetVerseRange(req.GetBiblePath(), req.GetInput())
	if err != nil {
		return nil, err
	}

	return &pb.VerseRange{
		Book:       book,
		Chapter:    int32(chapter),
		StartVerse: int32(startVerse),
		EndVerse:   int32(endVerse),
		Verses:     verseTexts,
	}, nil
}
