package gapi

import (
	"io"
	"log"

	"github.com/steve-mir/diivix_backend/internal/app/authentication/pb"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/services"
)

func (s *Server) UserSuggestions(stream pb.UserAuth_UserSuggestionsServer) error {
	// log.Println("Began")
	ctx := stream.Context()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			req, err := stream.Recv()

			if err == io.EOF {
				// Client has closed the stream
				return nil
			}

			if err != nil {
				log.Println("Error 1", err)
				return err
			}
			query := req.GetQuery()

			// Fetch user suggestions from Redis
			usernames, err := services.FetchUserSuggestions(s.redisCache, ctx, query)
			if err != nil {
				log.Println("Error 2", err)
				return err
			}

			resp := &pb.UserResponse{Usernames: usernames}
			if err := stream.Send(resp); err != nil {
				log.Println("Error 3", err)
				return err
			}
		}
	}
}
