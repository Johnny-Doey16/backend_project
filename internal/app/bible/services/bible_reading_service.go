package services

import (
	"fmt"
	"time"

	"github.com/steve-mir/diivix_backend/internal/app/bible/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const minDuration = 255

// * 8 months = 255 days
// * 1 year = 361 days
// * 1 year, 3 months = 476 days
// * 2 years = 686 days
// * 3 years, 6 months = 1342 days
func GenerateBibleReadingPlan(path string, durationInDays int, stream pb.BibleService_GenerateReadingPlanServer) error {

	if durationInDays < minDuration {
		return status.Errorf(codes.InvalidArgument, "Duration should be at least %d days", minDuration)
	}

	bible, err := loadTranslation(path)
	if err != nil {
		return status.Errorf(codes.Internal, "Error loading default translation: %s", err)
	}

	// Calculate the total number of chapters in the Bible
	totalChapters := 0
	for _, book := range bible.BibleBooks {
		totalChapters += len(book.Chapters)
	}

	// Calculate the number of readings per day
	readingsPerDay := totalChapters / durationInDays

	// Initialize variables
	currentDay := 1
	currentBookIndex := 0
	currentChapterIndex := 0
	currentVerseIndex := 0
	// scriptures := make([]string, 0)

	// Loop through the days
	for currentDay <= durationInDays {
		// Find the current book
		currentBook := &bible.BibleBooks[currentBookIndex]

		// Find the current chapter
		currentChapter := &currentBook.Chapters[currentChapterIndex]

		// Generate the scripture entry for the day
		startReading := fmt.Sprintf("%s %d:%d, ", currentBook.BName, currentChapter.CNumber, currentVerseIndex+1)

		// Calculate the end chapter for the day
		endChapter := currentChapter.CNumber + readingsPerDay - 1
		if endChapter >= len(currentBook.Chapters) {
			endChapter = len(currentBook.Chapters) - 1
		}

		// Append the end chapter to the scripture entry
		endReading := fmt.Sprintf("%s %d:%d", currentBook.BName, endChapter+1, len(currentBook.Chapters[endChapter].Verses))

		// Append the scripture entry to the result
		// scriptures = append(scriptures, scriptureEntry)
		biblePlan := &pb.BiblePlan{
			Day:          int32(currentDay),
			StartReading: startReading,
			EndReading:   endReading,
		}

		if err := stream.Send(biblePlan); err != nil {
			return status.Errorf(codes.Internal, "error sending verses. Error: %s", err)
		}

		// Update the indices for the next day
		currentVerseIndex = 0
		currentChapterIndex = endChapter + 1
		if currentChapterIndex >= len(currentBook.Chapters) {
			currentChapterIndex = 0
			currentBookIndex++

			// If all books are covered, reset to the first book
			if currentBookIndex >= len(bible.BibleBooks) {
				currentBookIndex = 0
			}
		}

		// Move to the next day
		currentDay++
	}

	// return scriptures, nil
	return nil
}

// Save plan to db. Start date(today), end date, days, bible(dra)
func SaveBibleReadingPlan(path string, durationInDays int) {
	startDate := time.Now()
	endDate := startDate.Add(time.Duration(durationInDays))
	bible := path
	fmt.Printf("Saving to db.\nPlan: %d, Started: %v, Ending: %v. Bible: %s", durationInDays, startDate, endDate, bible)
}

// GetCurrentReading returns the scripture reading for a specific day in the reading plan
func GetCurrentReading(path string, startDate time.Time, durationInDays int, currentDate time.Time) (*pb.BiblePlan, error) {
	// Calculate the number of days between the start date and the current date
	daysSinceStart := int(currentDate.Sub(startDate).Hours() / 24)

	// Ensure the day is within the planned duration
	if daysSinceStart < 1 || daysSinceStart > durationInDays {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid day requested")
	}

	// Calculate the day's reading based on the reading plan
	dayIndex := daysSinceStart - 1 // Adjust to zero-based index
	bible, err := loadTranslation(path)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error loading default translation: %s", err)
	}

	// Calculate the total number of chapters in the Bible
	totalChapters := 0
	for _, book := range bible.BibleBooks {
		totalChapters += len(book.Chapters)
	}

	// Calculate the number of readings per day
	readingsPerDay := totalChapters / durationInDays

	// Find the current book, chapter, and verse for the requested day
	currentBookIndex := 0
	currentChapterIndex := 0
	currentVerseIndex := 0

	for i := 1; i <= dayIndex; i++ {
		currentBook := &bible.BibleBooks[currentBookIndex]
		currentChapter := &currentBook.Chapters[currentChapterIndex]

		endChapter := currentChapter.CNumber + readingsPerDay - 1
		if endChapter >= len(currentBook.Chapters) {
			endChapter = len(currentBook.Chapters) - 1
		}

		currentVerseIndex = 0
		currentChapterIndex = endChapter + 1
		if currentChapterIndex >= len(currentBook.Chapters) {
			currentChapterIndex = 0
			currentBookIndex++

			// If all books are covered, reset to the first book
			if currentBookIndex >= len(bible.BibleBooks) {
				currentBookIndex = 0
			}
		}
	}

	// Generate the scripture entry for the requested day
	currentBook := &bible.BibleBooks[currentBookIndex]
	currentChapter := &currentBook.Chapters[currentChapterIndex]
	startReading := fmt.Sprintf("%s %d:%d, ", currentBook.BName, currentChapter.CNumber, currentVerseIndex+1)

	endChapter := currentChapter.CNumber + readingsPerDay - 1
	if endChapter >= len(currentBook.Chapters) {
		endChapter = len(currentBook.Chapters) - 1
	}

	endReading := fmt.Sprintf("%s %d:%d", currentBook.BName, endChapter+1, len(currentBook.Chapters[endChapter].Verses))

	return &pb.BiblePlan{
		Day:          int32(dayIndex + 1),
		StartReading: startReading,
		EndReading:   endReading,
	}, nil
}
