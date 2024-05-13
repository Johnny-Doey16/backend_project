package services

import (
	"strings"

	"github.com/steve-mir/diivix_backend/internal/app/bible/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// *
func SearchByKeyword(path, keyword string, stream pb.BibleService_SearchByKeywordServer) error {
	bible, err := loadTranslation(path)
	if err != nil {
		return status.Errorf(codes.Internal, "Error loading default translation: %s", err)
	}

	for _, book := range bible.BibleBooks {
		for _, chapter := range book.Chapters {
			for _, verse := range chapter.Verses {
				if strings.Contains(verse.Text, keyword) {

					resVerse := &pb.Verse{
						Book:    book.BName,
						Chapter: int32(chapter.CNumber),
						Verse:   int32(verse.VNumber),
						Text:    verse.Text,
					}

					if err := stream.Send(resVerse); err != nil {
						return status.Errorf(codes.Internal, "error sending verses. Error: %s", err)
					}
				}
			}
		}
	}
	return nil
}

func SearchByTopic(path, topic string, stream pb.BibleService_SearchByTopicServer) error {
	bible, err := loadTranslation(path)
	if err != nil {
		return status.Errorf(codes.Internal, "Error loading default translation: %s", err)
	}

	for _, book := range bible.BibleBooks {
		for _, chapter := range book.Chapters {
			for _, verse := range chapter.Verses {
				if strings.Contains(verse.Text, topic) {

					resVerse := &pb.Verse{
						Book:    book.BName,
						Chapter: int32(chapter.CNumber),
						Verse:   int32(verse.VNumber),
						Text:    verse.Text,
					}

					if err := stream.Send(resVerse); err != nil {
						return status.Errorf(codes.Internal, "error sending verses. Error: %s", err)
					}

				}
			}
		}
	}
	return nil
}
