package services

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// *
func SearchBibleVerse(path, input string) (string, int, int, string, error) {
	defaultTranslation, err := loadTranslation(path)
	if err != nil {
		return "", 0, 0, "", status.Errorf(codes.Internal, "Error loading default translation: %s", err)
	}

	return SearchVerse(defaultTranslation, input)
}

func SearchVerse(bible Bible, input string) (string, int, int, string, error) {
	book, chapter, verse, err := parseInput(input)
	if err != nil {
		return "", 0, 0, "", err
	}

	targetBook := findBook(bible.BibleBooks, book)
	if targetBook == nil {
		return "", 0, 0, "", status.Errorf(codes.NotFound, "book not found: %s", book)
	}

	targetChapter := findChapter(targetBook.Chapters, chapter)
	if targetChapter == nil {
		return "", 0, 0, "", status.Errorf(codes.NotFound, "chapter not found: %d", chapter)
	}

	targetVerse := findVerse(targetChapter.Verses, verse)
	if targetVerse == nil {
		return "", 0, 0, "", status.Errorf(codes.NotFound, "verse not found: %d", verse)
	}

	// return fmt.Sprintf("%s %d:%d - %s", book, chapter, verse, targetVerse.Text), nil
	return book, chapter, verse, targetVerse.Text, nil
}
