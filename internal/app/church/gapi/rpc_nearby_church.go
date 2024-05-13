package gapi

import (
	"context"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
)

const (
	// Define the search radius in meters (example: 5000 meters)
	radius = float64(5000) // FIXME: Test in the UI to Adjust
)

func (s *ChurchServer) SearchNearbyChurches(req *pb.SearchRequestChurch, stream pb.ChurchService_SearchNearbyChurchesServer) error {
	// Check if user_location is provided in the request
	if req.UserLocation == nil {
		// Handle the case where user location is not provided if required
		return nil // or an appropriate error
	}

	// Extract latitude and longitude from the request
	lat := req.UserLocation.Lat
	lng := req.UserLocation.Lng

	// Get the nearby churches using the generated sqlc method
	churches, err := s.store.GetNearbyChurches(context.Background(), sqlc.GetNearbyChurchesParams{
		StMakepoint:   lat,
		StMakepoint_2: lng,
		StDwithin:     radius,
	})
	if err != nil {
		// Handle the error properly
		return err
	}

	// Stream back the result to the client
	for _, church := range churches {
		// Convert the sqlc church type to protobuf church message
		// Assuming you have a conversion function or method to convert the types
		churchMsg := convertToChurchMessage(church)
		if err := stream.Send(churchMsg); err != nil {
			// Handle the error properly
			return err
		}
	}

	return nil
}

// Convert the sqlc generated church type to protobuf church message
// Implement this function based on your actual types
func convertToChurchMessage(church sqlc.GetNearbyChurchesRow) *pb.Church {
	// Conversion logic here
	return &pb.Church{
		AuthId:         church.AuthID.String(),
		Id:             church.ID,
		Name:           church.Name,
		ImageUrl:       church.ImageUrl.String,
		Username:       church.Username.String,
		DenominationId: church.DenominationID,
		Location: &pb.Location{
			Country:    church.Country,
			State:      church.State,
			City:       church.City,
			Address:    church.Address,
			PostalCode: church.Postalcode,
			Lga:        church.Lga,
		},
		// Email: church.Email,
		// Phone: church.Phone,
	}
}
