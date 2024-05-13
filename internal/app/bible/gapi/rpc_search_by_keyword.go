package gapi

import (
	"github.com/steve-mir/diivix_backend/internal/app/bible/pb"
	"github.com/steve-mir/diivix_backend/internal/app/bible/services"
)

func (s *BibleServer) SearchByKeyword(req *pb.SearchByKeywordRequest, stream pb.BibleService_SearchByKeywordServer) error {
	err := services.SearchByKeyword(req.GetBiblePath(), req.GetKeyword(), stream)
	if err != nil {
		return err
	}

	return nil
}
