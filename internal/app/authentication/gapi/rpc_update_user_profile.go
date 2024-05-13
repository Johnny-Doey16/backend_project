package gapi

import (
	"context"
	"time"

	"github.com/steve-mir/diivix_backend/internal/app/authentication/pb"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/services"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//  https://images.unsplash.com/photo-1701850975931-f3a971bf002d?q=80&w=3687&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D

func (s *Server) UpdateUserProfile(stream pb.UserAuth_UpdateUserProfileServer) error {
	ctx, cancel := context.WithTimeout(stream.Context(), time.Second*10) // Adjust the timeout as needed
	defer cancel()

	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	imageUrl, firstName, lastName, about, website, headerImageUrl, err := services.GeUserProfileStreamData(stream, s.imageStore)
	if err != nil {
		return err
	}

	err = services.CreateUserProfileInDB(ctx, s.db, s.store, claims.UserId, firstName, lastName, about, website, imageUrl, headerImageUrl)
	if err != nil {
		return err
	}

	// Send the response back to the client
	return stream.SendAndClose(&pb.UpdateUserProfileResponse{
		ImageUrl:       imageUrl,
		HeaderImageUrl: headerImageUrl,
		Status:         true,
		Msg:            "Profile updated successfully",
	})
}
