package services

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ListAllChaptersOfBook(path, bookName string) ([]int, error) {
	bible, err := loadTranslation(path)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error loading default translation: %s", err)
	}

	targetBook := findBook(bible.BibleBooks, bookName)
	if targetBook == nil {
		return nil, status.Errorf(codes.NotFound, "Book not found: "+bookName)
	}
	chapters := []int{}
	for _, chapter := range targetBook.Chapters {
		fmt.Println(chapter.CNumber)
		chapters = append(chapters, chapter.CNumber)
	}
	return chapters, nil
}
