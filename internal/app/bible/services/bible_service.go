package services

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// I already have the below code for some functions. Add another function for generating scriptures for the user. The function is for one who wants to read the whole Bible
// within a period of time (1 year, 6 months, 2 years etc). Assuming the user decides to read the entire Bible in a Year the function with split the all books and chapter for
// user in 1 year. They will now have readings for each day. The output would be something like this
// [Day: 1, Start: Gen 1:1, End: Gen 4:5, Day: 2, Start Gen 4:6, End: Gen 10:4...]. Note that the verses and chapters are not consistent so as not to generate a reading that is not
// in the Bible.

// Bible represents the structure of your XML data
type Bible struct {
	XMLName    xml.Name    `xml:"XMLBIBLE"`
	BibleBooks []BibleBook `xml:"BIBLEBOOK"`
}

// BibleBook represents a book in the Bible
type BibleBook struct {
	BNumber  int       `xml:"bnumber,attr"`
	BName    string    `xml:"bname,attr"`
	BSName   string    `xml:"bsname,attr"`
	Chapters []Chapter `xml:"CHAPTER"`
}

// Chapter represents a chapter in a book
type Chapter struct {
	CNumber int     `xml:"cnumber,attr"`
	Verses  []Verse `xml:"VERS"`
}

// Verse represents a verse in a chapter
type Verse struct {
	VNumber int    `xml:"vnumber,attr"`
	Text    string `xml:",chardata"`
}

func parseInput(input string) (string, int, int, error) {
	parts := strings.Split(input, " ")
	if len(parts) != 2 {
		return "", 0, 0, status.Errorf(codes.InvalidArgument, "Invalid format: use 'Book Chapter:Verse'")
	}

	book := parts[0]
	chapterVerseParts := strings.Split(parts[1], ":")
	if len(chapterVerseParts) != 2 {
		return "", 0, 0, status.Errorf(codes.InvalidArgument, "Invalid format: use 'Book Chapter:Verse'")
	}

	chapter, err := strconv.Atoi(chapterVerseParts[0])
	if err != nil {
		return "", 0, 0, status.Errorf(codes.InvalidArgument, "Invalid chapter number: %s", chapterVerseParts[0])
	}

	verse, err := strconv.Atoi(chapterVerseParts[1])
	if err != nil {
		return "", 0, 0, status.Errorf(codes.InvalidArgument, "Invalid verse number: %s", chapterVerseParts[1])
	}

	return book, chapter, verse, nil
}

func findBook(books []BibleBook, name string) *BibleBook {
	for _, book := range books {
		if book.BName == name {
			return &book
		}
	}
	return nil
}

func findChapter(chapters []Chapter, number int) *Chapter {
	for _, chapter := range chapters {
		if chapter.CNumber == number {
			return &chapter
		}
	}
	return nil
}

func findVerse(verses []Verse, number int) *Verse {
	for _, verse := range verses {
		if verse.VNumber == number {
			return &verse
		}
	}
	return nil
}

func findVerseRange(bookName string, chapter int, verses []Verse, startVerse, endVerse int) []string {
	// Check whether to use []Verse
	words := []string{}

	for _, verse := range verses {
		if verse.VNumber >= startVerse && verse.VNumber <= endVerse {
			word := fmt.Sprintf("%s %d:%d - %s\n", bookName, chapter, verse.VNumber, verse.Text)
			words = append(words, word)
			// fmt.Printf("%s %d:%d - %s\n", bookName, targetChapter.CNumber, verse.VNumber, verse.Text)
		}
	}

	return words
}

func loadTranslation(filename string) (Bible, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Bible{}, err
	}
	defer file.Close()

	var bible Bible
	if err := xml.NewDecoder(file).Decode(&bible); err != nil {
		return Bible{}, err
	}
	return bible, nil
}

// Helper function to get the text of a specific verse

// func getVerse(bible Bible, bookName string, chapterNumber, verseNumber int) string {
// 	targetBook := findBook(bible.BibleBooks, bookName)
// 	if targetBook == nil {
// 		return "" //fmt.Errorf("Book not found: %s", book)
// 	}

// 	targetChapter := findChapter(targetBook.Chapters, chapterNumber)
// 	if targetChapter == nil {
// 		return "" //fmt.Errorf("Chapter not found: %d", chapter)
// 	}

// 	targetVerse := findVerse(targetChapter.Verses, verseNumber)
// 	if targetVerse == nil {
// 		return "" // fmt.Errorf("Verse not found: %d", verse)
// 	}

// 	return targetVerse.Text
// }

// Helper function to get cross-references for a given verse
// func getCrossReferences(bookName string, chapterNumber, verseNumber int) []string {
// 	return []string{"Cross"}
// }
