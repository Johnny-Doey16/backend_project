package services

// import (
// 	"bytes"
// 	"context"
// 	"database/sql"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"log"
// 	"sync"
// 	"time"

// 	"github.com/google/uuid"
// 	"github.com/steve-mir/diivix_backend/cache"
// 	"github.com/steve-mir/diivix_backend/db/sqlc"
// 	"github.com/steve-mir/diivix_backend/internal/app/authentication/services"
// 	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )

// const (
// 	maxImages    = 4
// 	maxImageSize = 1 << 20
// 	minContent   = 20
// 	maxContent   = 280 //400
// )

// type postStruct struct {
// 	p   sqlc.Post
// 	err error
// }

// func GetCreatePostStreamData(stream pb.SocialMedia_CreatePostServer, imgStore services.ImageStore) ([]string, [maxImages]string, string, error) {
// 	var (
// 		imageUrls []string
// 		ext       string
// 		content   string
// 		imageData bytes.Buffer //imageData := bytes.Buffer{}
// 		imageSize int          // 0
// 		captions  [maxImages]string
// 		index     int // = 0 // used to determine  which index to add the caption to
// 	)

// 	for {
// 		req, err := stream.Recv()
// 		if err == io.EOF {
// 			// End of stream, process the data
// 			break
// 		}
// 		if err != nil {
// 			return []string{}, [maxImages]string{}, "", status.Errorf(codes.Internal, "error receiving stream from client: %v", err)
// 		}

// 		if err := checkBody(req); err != nil {
// 			return []string{}, [maxImages]string{}, "", status.Errorf(codes.InvalidArgument, err.Error())
// 		}

// 		switch data := req.Data.(type) {
// 		// has content
// 		case *pb.CreatePostRequest_Content:
// 			content = data.Content

// 			// has metadata
// 		case *pb.CreatePostRequest_Metadata:
// 			ext = data.Metadata.Ext
// 			captions[index] = data.Metadata.Caption

// 			// contains image
// 		case *pb.CreatePostRequest_ChunkData:
// 			imageErr := handleImageUpload(data, imgStore, ext, &imageUrls, &imageData, &imageSize, &index)
// 			if imageErr != nil {
// 				return []string{}, [maxImages]string{}, "", imageErr
// 			}

// 		}
// 	}

// 	return imageUrls, captions, content, nil

// }

// func checkBody(req *pb.CreatePostRequest) error {
// 	// ? Fix, does well with only text but when sent with image returns the error
// 	if req.GetContent() == "" && req.GetChunkData() == nil {
// 		return errors.New("body cannot be empty")
// 	}

// 	if len(req.GetContent()) < minContent {
// 		return fmt.Errorf("characters must be at least %v characters long", minContent)
// 	}

// 	if len(req.GetContent()) > maxContent {
// 		return fmt.Errorf("characters must not be more than %v characters long", maxContent)
// 	}
// 	return nil
// }

// func handleImageUpload(data *pb.CreatePostRequest_ChunkData, imageStore services.ImageStore, ext string, imageUrls *[]string, imageData *bytes.Buffer, imageSize *int, index *int) error {
// 	uid, err := uuid.NewRandom()
// 	if err != nil {
// 		return status.Errorf(codes.Internal, "error generating image UUID: %v", err)
// 	}

// 	imgURL := fmt.Sprintf("%s%s", uid.String(), ext)

// 	size := len(data.ChunkData)

// 	*imageSize += size
// 	if *imageSize > maxImageSize {
// 		log.Println("Image too large")
// 		// return status.Errorf(codes.Unknown, "image is too large: "+strconv.Itoa(*imageSize)+"> "+strconv.Itoa(maxImageSize))
// 	}

// 	if _, err := imageData.Write(data.ChunkData); err != nil {
// 		return status.Errorf(codes.Internal, "cannot write chunk data: %v", err)
// 	}

// 	if err := imageStore.Save(uid, ext, *imageData); err != nil {
// 		return status.Errorf(codes.Internal, "cannot save image to the store: %v", err)
// 	}

// 	*imageUrls = append(*imageUrls, imgURL)
// 	*index++

// 	// Reset the imageData buffer for the next image
// 	imageData.Reset()

// 	return nil
// }

// // Creates post in db with the content, image urls, caption etc
// func CreatePostInDB(ctx context.Context, db *sql.DB, store *sqlc.Store, uid uuid.UUID, content string, imageUrls []string, captions [maxImages]string) (uuid.UUID, sqlc.Post, error) {
// 	postID, err := uuid.NewRandom()
// 	if err != nil {
// 		return uuid.Nil, sqlc.Post{}, status.Errorf(codes.Internal, "error generating post ID: %v", err)
// 	}

// 	// Create transaction
// 	tx, err := db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return uuid.Nil, sqlc.Post{}, status.Errorf(codes.Internal, "transaction error: %v", err.Error())
// 	}
// 	defer tx.Rollback()
// 	qtx := store.WithTx(tx)

// 	sqlPost, err := runConcurrentPostCreationTasks(tx, qtx, ctx, postID, imageUrls, captions, content, uid)
// 	if err != nil {
// 		return uuid.Nil, sqlc.Post{}, status.Errorf(codes.Internal, "run concurrency error: %v", err)
// 	}

// 	if err := tx.Commit(); err != nil {
// 		return uuid.Nil, sqlc.Post{}, status.Errorf(codes.Internal, "commit: %v", err)
// 	}

// 	return postID, sqlPost, nil
// }

// func runConcurrentPostCreationTasks(tx *sql.Tx, qtx *sqlc.Queries, ctx context.Context,
// 	postId uuid.UUID, images []string, captions [maxImages]string, p_content string,
// 	uid uuid.UUID) (sqlc.Post, error) {

// 	var wg sync.WaitGroup
// 	createMetricChan := make(chan error, 1)
// 	createImgChan := make(chan error, 1)
// 	// createMentionsChan := make(chan error, 1)
// 	createPostChan := make(chan postStruct, 1)

// 	// ! 1 Start the post creation goroutine
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		sqlPost, err := createPost(qtx, ctx, postId, p_content, uid, len(images))
// 		createPostChan <- postStruct{p: sqlPost, err: err}
// 	}()

// 	// Wait for the post creation to complete before starting other operations
// 	wg.Wait()
// 	postData := <-createPostChan
// 	if postData.err != nil {
// 		tx.Rollback()
// 		return sqlc.Post{}, errors.New("error creating post " + postData.err.Error())
// 	}

// 	// ! 1b If post creation is successful, proceed with images and metrics
// 	if len(images) > 0 {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			i := 0
// 			for _, imgUrl := range images {
// 				err := qtx.CreatePostImage(ctx, sqlc.CreatePostImageParams{
// 					PostID:   postId,
// 					ImageUrl: sql.NullString{String: imgUrl, Valid: true},
// 					Caption:  sql.NullString{String: captions[i], Valid: true},
// 				})
// 				if err != nil {
// 					createImgChan <- err
// 					return
// 				}
// 				i++
// 			}
// 			createImgChan <- nil
// 		}()
// 	} else {
// 		createImgChan <- nil
// 	}

// 	// ! 2 create posts metrics
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		err := qtx.CreateMetric(ctx, sqlc.CreateMetricParams{
// 			PostID:   postId,
// 			Views:    sql.NullInt32{Int32: 0, Valid: true},
// 			Likes:    sql.NullInt32{Int32: 0, Valid: true},
// 			Reposts:  sql.NullInt32{Int32: 0, Valid: true},
// 			Comments: sql.NullInt32{Int32: 0, Valid: true},
// 		})
// 		createMetricChan <- err
// 	}()

// 	// ! 3 create mentions if necessary
// 	// wg.Add(1)
// 	// go func() {
// 	// 	defer wg.Done()
// 	// 	err := DetectMentions(ctx, p_content, tx, postId)
// 	// 	createMentionsChan <- err
// 	// }()
// 	//

// 	wg.Wait()
// 	close(createMetricChan)
// 	close(createImgChan)
// 	// close(createMentionsChan)

// 	if err := <-createImgChan; err != nil {
// 		tx.Rollback()
// 		return sqlc.Post{}, errors.New("an unknown error occurred creating image " + err.Error())
// 	}

// 	if err := <-createMetricChan; err != nil {
// 		tx.Rollback()
// 		return sqlc.Post{}, errors.New("an unknown error occurred creating metric " + err.Error())
// 	}

// 	// if err := <-createMentionsChan; err != nil {
// 	// 	// tx.Rollback()
// 	// 	return sqlc.Post{}, errors.New("an unknown error occurred creating mentions " + err.Error())
// 	// }

// 	return postData.p, nil
// }

// func createPost(store *sqlc.Queries, ctx context.Context, postId uuid.UUID, p_content string, uid uuid.UUID, imgLength int) (sqlc.Post, error) {

// 	return store.CreatePost(ctx, sqlc.CreatePostParams{
// 		ID:          postId,
// 		Content:     p_content,
// 		UserID:      uid,
// 		TotalImages: sql.NullInt32{Int32: int32(imgLength), Valid: true},
// 	})
// }

// func CachePost(ctx context.Context, redisCache cache.Cache, post interface{} /*post *pb.Post*/, key string, expr time.Duration) error {
// 	// Use the post ID as the key and serialize the post object as the value.
// 	// You can use JSON or any other serialization method of your choice.
// 	// Here's an example using JSON:

// 	serializedPost, err := json.Marshal(post)
// 	if err != nil {
// 		return err
// 	}

// 	// Cache the serialized post in Redis with an expiration time (e.g., 24 hours).
// 	// The key pattern here is "post:{postID}". fmt.Sprintf("posts:%s", post.PostId)
// 	return redisCache.SetKey(ctx, key, serializedPost, expr) //.Err()
// }
