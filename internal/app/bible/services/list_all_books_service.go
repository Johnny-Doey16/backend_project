package services

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ListAllBooksOfBible(path string) ([]string, error) {
	bible, err := loadTranslation(path)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error loading default translation: %s", err)
	}

	books := []string{}
	for _, book := range bible.BibleBooks {
		books = append(books, book.BName)
	}
	return books, nil
}
