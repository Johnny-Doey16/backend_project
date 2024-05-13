package services

import (
	"fmt"

	"github.com/steve-mir/diivix_backend/internal/app/bible/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ListAllVersesOfChapter(path, bookName string, chapterNumber int, stream pb.BibleService_ListAllVersesOfChapterServer) error {
	bible, err := loadTranslation(path)
	if err != nil {
		return status.Errorf(codes.Internal, "Error loading default translation: %s", err)
	}

	targetBook := findBook(bible.BibleBooks, bookName)
	if targetBook == nil {
		return status.Errorf(codes.NotFound, "Book not found: %s", bookName)
	}

	targetChapter := findChapter(targetBook.Chapters, chapterNumber)
	if targetChapter == nil {
		return status.Errorf(codes.NotFound, "Chapter not found: %d", chapterNumber)
	}

	fmt.Printf("List of all verses in %s %d:\n", bookName, chapterNumber)
	for _, verse := range targetChapter.Verses {

		if err := stream.Send(&pb.Verse{
			Book:    bookName,
			Chapter: int32(chapterNumber),
			Verse:   int32(verse.VNumber),
			Text:    verse.Text,
		}); err != nil {
			return status.Errorf(codes.Internal, "error sending verses of %s chapter %d. Error: %s", bookName, chapterNumber, err)
		}
	}
	return nil
}
