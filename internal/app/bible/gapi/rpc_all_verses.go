package gapi

import (
	"github.com/steve-mir/diivix_backend/internal/app/bible/pb"
	"github.com/steve-mir/diivix_backend/internal/app/bible/services"
)

func (s *BibleServer) ListAllVersesOfChapter(req *pb.ListAllVersesOfChapterRequest, stream pb.BibleService_ListAllVersesOfChapterServer) error {
	return services.ListAllVersesOfChapter(req.GetBiblePath(), req.GetBookName(), int(req.GetChapterNumber()), stream)
}
