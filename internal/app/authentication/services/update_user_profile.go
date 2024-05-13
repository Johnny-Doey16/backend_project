package services

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxImageSize = 1 << 20

func GeUserProfileStreamData(stream pb.UserAuth_UpdateUserProfileServer, imgStore ImageStore) (string, string, string, string, string, string, error) {
	var (
		imageUrl        string
		ext             string
		firstName       string
		lastName        string
		imageData       bytes.Buffer //imageData := bytes.Buffer{}
		imageSize       int          // 0
		website         string
		about           string
		headerImageData bytes.Buffer
		headerImageUrl  string
		headerExt       string
	)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// End of stream, process the data
			break
		}
		if err != nil {
			return "", "", "", "", "", "", status.Errorf(codes.Internal, "error receiving stream from client: %v", err)
		}

		switch data := req.Data.(type) {
		// has content
		case *pb.UpdateUserProfileRequest_FirstName:
			firstName = data.FirstName

		case *pb.UpdateUserProfileRequest_LastName:
			lastName = data.LastName

			// has metadata
		case *pb.UpdateUserProfileRequest_ProfileImgMetadata:
			ext = data.ProfileImgMetadata.Ext

			// has HeaderMetadata
		case *pb.UpdateUserProfileRequest_HeaderImgMetadata:
			headerExt = data.HeaderImgMetadata.Ext

			// has about
		case *pb.UpdateUserProfileRequest_About:
			about = data.About

			// has website
		case *pb.UpdateUserProfileRequest_Website:
			website = data.Website

			// contains image
		case *pb.UpdateUserProfileRequest_ProfileImgData:
			imageErr := handleImageUpload(data, imgStore, ext, &imageUrl, &imageData, &imageSize)
			if imageErr != nil {
				return "", "", "", "", "", "", imageErr
			}

			// contains header image
		case *pb.UpdateUserProfileRequest_HeaderImgData:
			imageErr := handleHeaderImageUpload(data, imgStore, headerExt, &headerImageUrl, &headerImageData, &imageSize)
			if imageErr != nil {
				return "", "", "", "", "", "", imageErr
			}

		}
	}

	return imageUrl, firstName, lastName, about, website, headerImageUrl, nil

}

// Creates post in db with the content, image urls, caption etc
func CreateUserProfileInDB(ctx context.Context, db *sql.DB, store *sqlc.Store, uid uuid.UUID, firstName, lastName, about, website, imageUrl, headerImageUrl string) error {
	// Create transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return status.Errorf(codes.Internal, "transaction error: %v", err.Error())
	}
	defer tx.Rollback()
	qtx := store.WithTx(tx)

	err = runConcurrentProfileUpdateTasks(tx, qtx, ctx, uid, firstName, lastName, about, website, imageUrl, headerImageUrl)
	if err != nil {
		return status.Errorf(codes.Internal, "run concurrency error: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return status.Errorf(codes.Internal, "commit: %v", err)
	}

	return nil
}

func handleImageUpload(data *pb.UpdateUserProfileRequest_ProfileImgData, imageStore ImageStore, ext string, imageUrls *string, imageData *bytes.Buffer, imageSize *int) error {
	uid, err := uuid.NewRandom()
	if err != nil {
		return status.Errorf(codes.Internal, "error generating image UUID: %v", err)
	}

	imgURL := fmt.Sprintf("%s%s", uid.String(), ext)

	size := len(data.ProfileImgData)

	*imageSize += size
	if *imageSize > maxImageSize {
		// log.Println("Image too large")
		return status.Errorf(codes.Unknown, "image is too large: "+strconv.Itoa(*imageSize)+"> "+strconv.Itoa(maxImageSize))
	}

	if _, err := imageData.Write(data.ProfileImgData); err != nil {
		return status.Errorf(codes.Internal, "cannot write chunk data: %v", err)
	}

	// Upload to aws
	// err = awsservice.UploadToS3(server.config, compressedImage, UID.String(), ext)
	// if err != nil {
	// 	return status.Errorf(codes.Internal, "cannot upload image to aws: "+err.Error())
	// }
	if err := imageStore.Save(uid, ext, *imageData); err != nil {
		return status.Errorf(codes.Internal, "cannot save image to the store: %v", err)
	}

	*imageUrls = imgURL

	return nil
}

// TODO: Change to upload image to cloud (aws, firebase)
func handleHeaderImageUpload(data *pb.UpdateUserProfileRequest_HeaderImgData, imageStore ImageStore, ext string, imageUrls *string, imageData *bytes.Buffer, imageSize *int) error {
	uid, err := uuid.NewRandom()
	if err != nil {
		return status.Errorf(codes.Internal, "error generating image UUID: %v", err)
	}

	imgURL := fmt.Sprintf("%s%s", uid.String(), ext)

	size := len(data.HeaderImgData)

	*imageSize += size
	if *imageSize > maxImageSize {
		// log.Println("Image too large")
		return status.Errorf(codes.Unknown, "image is too large: "+strconv.Itoa(*imageSize)+"> "+strconv.Itoa(maxImageSize))
	}

	if _, err := imageData.Write(data.HeaderImgData); err != nil {
		return status.Errorf(codes.Internal, "cannot write chunk data: %v", err)
	}

	// Upload to aws
	// err = awsservice.UploadToS3(server.config, compressedImage, UID.String(), ext)
	// if err != nil {
	// 	return status.Errorf(codes.Internal, "cannot upload image to aws: "+err.Error())
	// }
	if err := imageStore.Save(uid, ext, *imageData); err != nil {
		return status.Errorf(codes.Internal, "cannot save image to the store: %v", err)
	}

	*imageUrls = imgURL

	return nil
}

func runConcurrentProfileUpdateTasks(tx *sql.Tx, qtx *sqlc.Queries, ctx context.Context,
	uid uuid.UUID, firstName, lastName, about, website, imgUrl, headerImageUrl string,
) error {

	var wg sync.WaitGroup
	updateUsersChan := make(chan error, 1) // Update user
	updateImgChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := qtx.UpdateUserNames(ctx, sqlc.UpdateUserNamesParams{
			UserID:    uid,
			FirstName: sql.NullString{String: firstName, Valid: firstName != ""},
			LastName:  sql.NullString{String: lastName, Valid: lastName != ""},
		})
		updateUsersChan <- err
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := qtx.UpdateImgEntityProfile(ctx, sqlc.UpdateImgEntityProfileParams{
			UserID:         uid,
			About:          sql.NullString{String: about, Valid: true},
			Website:        sql.NullString{String: website, Valid: true},
			HeaderImageUrl: sql.NullString{String: headerImageUrl, Valid: true},
			ImageUrl:       sql.NullString{String: imgUrl, Valid: imgUrl != ""},
		})
		updateImgChan <- err
	}()

	wg.Wait()
	close(updateUsersChan)
	close(updateImgChan)

	if err := <-updateUsersChan; err != nil {
		tx.Rollback()
		return errors.New("an unknown error occurred creating metric " + err.Error())
	}

	if err := <-updateImgChan; err != nil {
		tx.Rollback()
		return errors.New("an unknown error occurred creating mentions " + err.Error())
	}

	return nil
}
