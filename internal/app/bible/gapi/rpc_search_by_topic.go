package gapi

import (
	"github.com/steve-mir/diivix_backend/internal/app/bible/pb"
	"github.com/steve-mir/diivix_backend/internal/app/bible/services"
)

func (s *BibleServer) SearchByTopic(req *pb.SearchByTopicRequest, stream pb.BibleService_SearchByTopicServer) error {
	err := services.SearchByTopic(req.GetBiblePath(), req.GetTopic(), stream)
	if err != nil {
		return err
	}
	return nil
}
