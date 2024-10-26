package s3service

import (
	"fmt"
	"github.com/EularGauss/bandlab-assignment/internal/app"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"time"
)

// S3ServiceImpl is the implementation of S3Service
type S3ServiceImpl struct {
	Bucket string
}

// NewS3Service creates a new S3Service instance
func NewS3Service(bucket string) *S3ServiceImpl {
	return &S3ServiceImpl{
		Bucket: bucket,
	}
}

// GeneratePresignedURL generates a pre-signed URL for uploading an image to S3
func (s *S3ServiceImpl) GeneratePresignedURL(key string) (string, error) {
	s3_config := app.DefaultS3Config()
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s3_config.Region), // e.g., us-west-2
	})
	if err != nil {
		return "", err
	}

	svc := s3.New(sess)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if req == nil {
		return "", fmt.Errorf("failed to create a new PutObjectRequest")
	}

	// Set expiration for the URL
	url, err := req.Presign(time.Duration(s3_config.PreSignedURLExpiration) * time.Minute) // URL is valid for 15 minutes
	if err != nil {
		return "", err
	}

	return url, nil
}
