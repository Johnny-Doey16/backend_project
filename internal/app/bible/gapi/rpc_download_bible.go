package gapi

import (
	"io"
	"os"

	"github.com/steve-mir/diivix_backend/internal/app/bible/pb"
)

func (s *BibleServer) DownloadBible(req *pb.DownloadBibleRequest, stream pb.BibleService_DownloadBibleServer) error {
	var path string

	switch req.BibleTranslation {
	case "dra":
		path = "bible-translations/eng-dra.osis.xml"
	case "kjv":
		path = "bible-translations/eng-kjv.xml"
	default:
		path = "bible-translations/eng-dra.osis.xml"
	}

	// Open the file to be streamed
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a buffer to read and send file data in chunks
	buffer := make([]byte, 1024)

	// Loop to read file and send it in chunks
	for {
		// Read a chunk of data from the file
		bytesRead, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				// End of file reached, break the loop
				break
			}
			return err
		}

		// Send the chunk of data to the client
		if err := stream.Send(&pb.DownloadBibleResponse{Data: buffer[:bytesRead]}); err != nil {
			return err
		}
	}

	return nil
}
