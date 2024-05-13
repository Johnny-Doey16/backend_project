package gapi

import (
	"bytes"
	"context"
	"io"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/pb"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"

	// ut "github.com/steve-mir/diivix_backend/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ProfileServer) UploadImage(stream pb.UserProfile_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot receive image info")
	}

	exif := req.GetMetadata().GetExif()
	imageExt := req.GetMetadata().GetExt()
	log.Printf("receive an upload-image request for image %s", imageExt)

	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		log.Print("waiting to receive more data")

		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot receive chunk data: "+err.Error())
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		imageSize += size
		if imageSize > maxImageSize {
			return status.Errorf(codes.Unknown, "image is too large: "+strconv.Itoa(imageSize)+"> "+strconv.Itoa(maxImageSize))
		}

		_, err = imageData.Write(chunk)
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot write chunk data: "+err.Error())
		}
	}

	UID, err := uuid.NewRandom()
	if err != nil {
		return status.Errorf(codes.Internal, "cannot generate id: "+err.Error())
	}

	ext := req.GetMetadata().Ext

	// Upload to aws
	// err = awsservice.UploadToS3(server.config, compressedImage, UID.String(), ext)
	// if err != nil {
	// 	return status.Errorf(codes.Internal, "cannot upload image to aws: "+err.Error())
	// }

	// Save image to file
	err = server.imageStore.Save(UID, ext, imageData)
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot save image to the store: "+err.Error())
	}

	res := &pb.UploadImageResponse{
		Urls: []string{UID.String() + ext},
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot send response: "+err.Error())
	}

	log.Printf("saved image with id: %s, size: %d, exif: %s", UID, imageSize, exif)
	return nil
}

// Update user profile info
func (server *ProfileServer) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	_, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	// Send data to db
	// err := server.store.UpdateUserProfile(ctx, sqlc.UpdateUserProfileParams{
	// 	UserID:    claims.UserId,
	// 	FirstName: sql.NullString{String: req.GetFirstName(), Valid: req.FirstName != nil},
	// 	LastName:  sql.NullString{String: req.GetLastName(), Valid: req.LastName != nil},
	// 	ImageUrl:  sql.NullString{String: req.GetImageUrl(), Valid: req.ImageUrl != nil},
	// 	UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	// })
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "error updating profile %s", err)
	// }

	return &pb.UpdateProfileResponse{
		// Uid:    req.GetImageUrl(),
		Msg: "Profile updated successfully",
	}, nil
}
