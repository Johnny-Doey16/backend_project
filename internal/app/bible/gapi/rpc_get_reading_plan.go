package gapi

import (
	"context"

	"github.com/steve-mir/diivix_backend/internal/app/bible/pb"
	"github.com/steve-mir/diivix_backend/internal/app/bible/services"
)

func (s *BibleServer) GetCurrentReading(ctx context.Context, req *pb.GetCurrentReadingRequest) (*pb.BiblePlan, error) {

	// Assuming the reading plan started on November 1st
	// startDate := time.Date(2023, time.November, 1, 0, 0, 0, 0, time.UTC)

	// // Current date is February 17th
	// currentDate := time.Date(2024, time.February, 17, 0, 0, 0, 0, time.UTC)

	// // Reading plan duration is 255 days
	// durationInDays := 255

	// // Path to the XML file containing the Bible data
	// biblePath := "../bible-translations/eng-dra.osis.xml"

	// Get the current reading for the specified day
	// currentReading, err := services.GetCurrentReading(req.GetBiblePath(), req.GetStartDate().AsTime(), int(req.GetDurationInDays()), req.CurrentDate.AsTime())
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return nil, err
	// }

	// fmt.Println("Current Reading:")
	// fmt.Println(currentReading)

	return services.GetCurrentReading(req.GetBiblePath(), req.GetStartDate().AsTime(), int(req.GetDurationInDays()), req.CurrentDate.AsTime())
}
