package services

import (
	"math/rand"
	"time"
)

// TODO: Handle error here and fix the random issue
func RandomVerseGenerator(path string, isToday bool) (string, int, int, string) {
	bible, err := loadTranslation(path)
	if err != nil {
		return "", 0, 0, err.Error()
	}

	// Seed the random number generator for true randomness
	if isToday {
		rand.Seed(int64(time.Now().Day()))
	} else {
		rand.Seed(time.Now().UnixNano())
	}

	// select a random book, chapter and verse
	randomBook := bible.BibleBooks[rand.Intn(len(bible.BibleBooks))]
	randomChapter := randomBook.Chapters[rand.Intn(len(randomBook.Chapters))]
	randomVerse := randomChapter.Verses[rand.Intn(len(randomChapter.Verses))]

	return randomBook.BName, randomChapter.CNumber, randomVerse.VNumber, randomVerse.Text
}
