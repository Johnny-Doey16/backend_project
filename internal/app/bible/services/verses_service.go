package services

import (
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// *
func GetVerseRange(path, input string) (string, int, int, int, []string, error) {
	bible, err := loadTranslation(path)
	if err != nil {
		return "", 0, 0, 0, nil, status.Errorf(codes.Internal, "Error loading default translation: %s", err)
	}

	// Split input into book name, chapter number, and verse range
	parts := strings.Split(input, " ")
	if len(parts) != 2 {
		return "", 0, 0, 0, nil, status.Errorf(codes.InvalidArgument, "Invalid input format. Please use 'Book Chapter:Verse-Verse'")
	}

	bookName := parts[0]

	// Split chapter and verse range into separate strings
	chapterVerseRangeParts := strings.Split(parts[1], ":")
	if len(chapterVerseRangeParts) != 2 {
		return "", 0, 0, 0, nil, status.Errorf(codes.InvalidArgument, "Invalid input format. Please use 'Book Chapter:Verse-Verse'")
	}

	// Attempt to convert chapter number and verse range
	chapterNumber, err := strconv.Atoi(chapterVerseRangeParts[0])
	if err != nil {
		return "", 0, 0, 0, nil, status.Errorf(codes.InvalidArgument, "Invalid chapter number: %s", chapterVerseRangeParts[0])
	}

	verseRangeParts := strings.Split(chapterVerseRangeParts[1], "-")
	if len(verseRangeParts) != 2 {
		return "", 0, 0, 0, nil, status.Errorf(codes.InvalidArgument, "Invalid verse range format: %s", chapterVerseRangeParts[1])
	}

	startVerse, err := strconv.Atoi(verseRangeParts[0])
	if err != nil {
		return "", 0, 0, 0, nil, status.Errorf(codes.InvalidArgument, "Invalid start verse: %s", verseRangeParts[0])
	}

	endVerse, err := strconv.Atoi(verseRangeParts[1])
	if err != nil {
		return "", 0, 0, 0, nil, status.Errorf(codes.InvalidArgument, "Invalid end verse: %s", verseRangeParts[1])
	}

	targetBook := findBook(bible.BibleBooks, bookName)
	if targetBook == nil {
		return "", 0, 0, 0, nil, status.Errorf(codes.NotFound, "Book not found: %s", bookName)
	}

	targetChapter := findChapter(targetBook.Chapters, chapterNumber)
	if targetChapter == nil {
		return "", 0, 0, 0, nil, status.Errorf(codes.NotFound, "Chapter not found in book: %s", bookName)
	}

	// Find the chapter and verses in the specified range
	words := findVerseRange(bookName, chapterNumber, targetChapter.Verses, startVerse, endVerse)
	return bookName, chapterNumber, startVerse, endVerse, words, nil
}
