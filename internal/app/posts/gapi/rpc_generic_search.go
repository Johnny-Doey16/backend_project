package gapi

import (
	"context"
	"database/sql"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// // Parse the user ID
// uid, err := services.StrToUUID("558c60aa-977f-4d38-885b-e813656371ac")
//
//	if err != nil {
//		return nil, status.Errorf(codes.Internal, "parsing uid: %s", err)
//	}
func (s *SocialMediaServer) GenericSearch(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResult, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	// Ensure that the page number and page size are positive.
	pageNumber := req.PageNumber
	pageSize := req.PageSize
	if pageNumber < 1 {
		pageNumber = 1
	}
	if pageSize < 1 {
		pageSize = 10 // Default page size to 10 if not specified or negative
	}

	// Calculate the offset based on the page number and page size.
	offset := (pageNumber - 1) * pageSize

	// Retrieve search results from the store

	// results, err := s.store.GetSearchResult(ctx, sqlc.GetSearchResultParams{
	// 	PlaintoTsquery: req.GetQuery(),
	// 	FollowerUserID: claims.UserId,
	// 	Limit:          pageSize,
	// 	Offset:         int32(offset),
	// })

	results, err := s.store.GetSearchResultOldNoFTS(ctx, sqlc.GetSearchResultOldNoFTSParams{
		Column1:        sql.NullString{String: req.GetQuery(), Valid: true},
		FollowerUserID: claims.UserId,
		Limit:          pageSize,
		Offset:         int32(offset),
	})
	if err != nil {
		return nil, err
	}

	// Create a search result object
	var searchResult pb.SearchResult

	// Populate the search result based on the source type
	for _, result := range results {
		switch result.Source {
		case "users":
			users, err := services.UnmarshalUserResult(result.Data)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "error fetching users: %s", err)
			}

			searchResult.UserData = append(searchResult.UserData, users...)

		case "churches":
			churches, err := services.UnmarshalChurchResult(result.Data)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "error fetching churches: %s", err)
			}
			searchResult.ChurchData = append(searchResult.ChurchData, churches...)

		case "posts":
			// log.Println("Raw JSON data:", string(result.Data))
			posts, err := services.UnmarshalPostResult(result.Data)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "error fetching posts: %s", err)
			}
			searchResult.PostData = append(searchResult.PostData, posts...)
		}
	}

	// Set HasMore based on the number of items in each category
	searchResult.HasMore = int32(len(searchResult.UserData)+len(searchResult.ChurchData)+len(searchResult.PostData)) >= pageSize

	return &searchResult, nil
}
