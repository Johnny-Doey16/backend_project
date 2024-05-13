package main

/*import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

type OSIS struct {
	XMLName  xml.Name `xml:"osis"`
	Text     string   `xml:",chardata"`
	Xmlns    string   `xml:"xmlns,attr"`
	Xsi      string   `xml:"xsi,attr"`
	Schema   string   `xml:"schemaLocation,attr"`
	OsisText struct {
		Text        string `xml:",chardata"`
		OsisIDWork  string `xml:"osisIDWork,attr"`
		OsisRefWork string `xml:"osisRefWork,attr"`
		XmlLang     string `xml:"lang,attr"`
		Canonical   string `xml:"canonical,attr"`
		Header      struct {
			Text         string         `xml:",chardata"`
			RevisionDesc []RevisionDesc `xml:"revisionDesc"`
			Work         Work           `xml:"work"`
		} `xml:"header"`
		Div struct {
			Text      string    `xml:",chardata"`
			Type      string    `xml:"type,attr"`
			Canonical string    `xml:"canonical,attr"`
			Title     string    `xml:"title"`
			Book      []KJVBook `xml:"div"`
		} `xml:"div"`
	} `xml:"osisText"`
}

type RevisionDesc struct {
	Text string `xml:",chardata"`
	Date string `xml:"date"`
	P    string `xml:"p"`
}

type Work struct {
	Text        string      `xml:",chardata"`
	OsisWork    string      `xml:"osisWork,attr"`
	Title       Title       `xml:"title"`
	Description Description `xml:"description"`
	Publisher   Publisher   `xml:"publisher"`
	Identifier  Identifier  `xml:"identifier"`
	Language    Language    `xml:"language"`
	Rights      Rights      `xml:"rights"`
	RefSystem   string      `xml:"refSystem"`
}

type Title struct {
	Text  string `xml:",chardata"`
	Type  string `xml:"type,attr"`
	Short string `xml:"short,attr"`
}

type Description struct {
	Text string `xml:",chardata"`
	Type string `xml:"type,attr"`
}

type Publisher struct {
	Text string `xml:",chardata"`
	Type string `xml:"type,attr"`
}

type Identifier struct {
	Text string `xml:",chardata"`
	Type string `xml:"type,attr"`
}

type Language struct {
	Text string `xml:",chardata"`
	Type string `xml:"type,attr"`
}

type Rights struct {
	Text string `xml:",chardata"`
	Type string `xml:"type,attr"`
}

type KJVBook struct {
	Text      string      `xml:",chardata"`
	Type      string      `xml:"type,attr"`
	OsisID    string      `xml:"osisID,attr"`
	Canonical string      `xml:"canonical,attr"`
	Title     Title       `xml:"title"`
	Chapter   []Chapter   `xml:"chapter"`
	P         []Paragraph `xml:"p"`
	Verse     []Verse     `xml:"verse"`
}

type Chapter struct {
	Text    string `xml:",chardata"`
	OsisRef string `xml:"osisRef,attr"`
	SID     string `xml:"sID,attr"`
	N       string `xml:"n,attr"`
}

type Paragraph struct {
	Text        string        `xml:",chardata"`
	Verse       []Verse       `xml:"verse"`
	TransChange []TransChange `xml:"transChange"`
}

type Verse struct {
	Text   string `xml:",chardata"`
	OsisID string `xml:"osisID,attr"`
	SID    string `xml:"sID,attr"`
	N      string `xml:"n,attr"`
}

type TransChange struct {
	Type string `xml:"type,attr"`
	Text string `xml:",chardata"`
}

func main() {

	xmlFile, err := os.Open("../bible-translations/eng-kjv.osis.xml")
	if err != nil {
		fmt.Println("Error opening XML file:", err)
		return
	}
	defer xmlFile.Close()

	var osis OSIS
	err = xml.NewDecoder(xmlFile).Decode(&osis)
	if err != nil {
		fmt.Println("Error decoding XML:", err)
		return
	}

	// ! KVJ
	// fmt.Println(osis.OsisText.Div.Book[0].Title)
	fmt.Println(osis.OsisText.Div.Book[0].P[0])
	// fmt.Println(osis.OsisText.Div.Book[0].Verse)
	// getBookKJV("Genesis", osis.OsisText.Div.Book)
	// fmt.Println(osis.OsisText.Div.Book[1].Chapter.N)
	// fmt.Println(osis.OsisText.Div.Book[0].OsisID) //.Title.Short)

	// for _, paragraph := range osis.OsisText.Div.Book[0].P {
	// 	fmt.Println(paragraph.Text)
	// 	for _, verse := range paragraph.Verse {
	// 		fmt.Println("Verse:", verse.Text)
	// 	}
	// 	for _, transChange := range paragraph.TransChange {
	// 		fmt.Printf("TransChange type=%s: %s\n", transChange.Type, transChange.Text)
	// 	}
	// }

}

func getBookKJV(book string, books []KJVBook) {
	for _, b := range books {
		if strings.EqualFold(b.Title.Short, book) {
			// return b
			fmt.Printf("Found: %s %d\n", b.OsisID, len(b.Chapter)/2) //b.Chapter.N)
			// fmt.Println(b.Chapter[0])
			// for _, verse := range b.Verse {
			// 	fmt.Printf("%s %s:%s %s\n", book, b.Chapter.N, verse.N, verse.Text)
			// }
			for _, p := range b.P {
				// fmt.Println(p)
				for _, v := range p.Verse {
					fmt.Println(v)
				}
				// if verse.OsisID == "Gen.1.10" {
				// 	fmt.Printf("%s %d:%s %s\n", book, len(b.Chapter)/2, verse.N, verse.Text)
				// 	return
				// }
			}
		}
	}
}
*/
