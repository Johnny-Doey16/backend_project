package awsservice

import (
	"bytes"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/steve-mir/diivix_backend/utils"
)

func UploadToS3(config utils.Config, image *bytes.Buffer, key, ext string) error {
	// Specify your AWS credentials and region directly in the aws.Config struct
	awsConfig := aws.Config{
		Region:      aws.String(config.AwsRegion), // replace with your AWS region
		Credentials: credentials.NewStaticCredentials(config.AwsS3AccessKey, config.AwsS3SecretKey, ""),
	}

	sess, err := session.NewSession(&awsConfig)
	if err != nil {
		return err
	}

	// Create an S3 service client
	svc := s3.New(sess)

	fileType := ext[1:]
	// Upload the file to the S3 bucket
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(config.AwsBucketName),
		Key:    aws.String(key + ext),
		// ACL:                  aws.String("public-read"), // replace with the desired ACL
		Body:                 bytes.NewReader(image.Bytes()),
		ContentLength:        aws.Int64(int64(image.Len())),
		ContentType:          aws.String(fmt.Sprintf("image/%s", fileType)), // replace with the appropriate content type
		ServerSideEncryption: aws.String("AES256"),                          // replace with the desired encryption method
	})
	if err != nil {
		return err
	}

	fmt.Printf("Successfully uploaded to %s\n", config.AwsBucketName)
	return nil
}

func GetDownloadURL(config utils.Config, key string) (string, error) {
	// Specify your AWS credentials and region directly in the aws.Config struct
	awsConfig := aws.Config{
		Region:      aws.String(config.AwsRegion), // replace with your AWS region
		Credentials: credentials.NewStaticCredentials(config.AwsS3AccessKey, config.AwsS3SecretKey, ""),
	}

	sess, err := session.NewSession(&awsConfig)
	if err != nil {
		return "", err
	}

	// Create an S3 service client
	svc := s3.New(sess)

	// Generate a pre-signed URL for the object with a specific expiration time (e.g., 1 hour)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(config.AwsBucketName),
		Key:    aws.String(key),
	})
	url, err := req.Presign(2 * time.Hour) // 1-hour expiration time
	if err != nil {
		return "", err
	}

	return url, nil
}
