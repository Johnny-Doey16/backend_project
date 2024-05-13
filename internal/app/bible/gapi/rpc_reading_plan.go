package gapi

import (
	"github.com/steve-mir/diivix_backend/internal/app/bible/pb"
	"github.com/steve-mir/diivix_backend/internal/app/bible/services"
)

func (s *BibleServer) GenerateReadingPlan(req *pb.GenerateReadingPlanRequest, stream pb.BibleService_GenerateReadingPlanServer) error {
	err := services.GenerateBibleReadingPlan(req.GetBiblePath(), int(req.GetDurationInDays()), stream)
	if err != nil {
		return err
	}

	// Save plan to db. Start date(today), end date, days, bible(dra)
	// Save to db
	services.SaveBibleReadingPlan(req.GetBiblePath(), int(req.GetDurationInDays()))
	return nil
}
