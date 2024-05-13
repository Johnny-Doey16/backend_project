package gapi

import (
	"bytes"
	"context"
	"io"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/pb"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/services"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxImageSize = 1 << 20

// @depreciated file depreciated
func (s *Server) UploadProfileImage(stream pb.UserAuth_UploadProfileImageServer) error {
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
	err = s.imageStore.Save(UID, ext, imageData)
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

func (s *Server) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}

	err := services.UpdateProfile(ctx, s.store, claims.UserId, req)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error updating profile %s", err)
	}

	return &pb.UpdateProfileResponse{
		Msg: "Profile updated successfully",
	}, nil
}
