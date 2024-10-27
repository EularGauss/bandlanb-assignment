package s3service

import (
	"fmt"
	"github.com/EularGauss/bandlab-assignment/internal/app"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3ServiceImpl represents the S3 service
type S3ServiceImpl struct {
	Bucket string
}

// GeneratePresignedURL generates a pre-signed URL for uploading an image to S3
func (s *S3ServiceImpl) GeneratePresignedURL(key string) (string, error) {
	// Check for allowed file extensions
	allowedExtensions := []string{".jpg", ".jpeg", ".bmp"}    // Added .jpeg for better compatibility
	ext := strings.ToLower(key[strings.LastIndex(key, "."):]) // Extract the file extension

	if !contains(allowedExtensions, ext) {
		return "", fmt.Errorf("unsupported file type: %s. Supported types are: %v", ext, allowedExtensions)
	}

	s3_config := app.DefaultS3Config()
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s3_config.Region), // e.g., us-west-2
	})
	if err != nil {
		return "", err
	}

	svc := s3.New(sess)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(key),
		ContentType: aws.String(getContentType(ext)), // Set ContentType based on the file extension
	})
	if req == nil {
		return "", fmt.Errorf("failed to create a new PutObjectRequest")
	}

	// Set expiration for the URL
	url, err := req.Presign(time.Duration(s3_config.PreSignedURLExpiration) * time.Minute) // URL is valid for x minutes
	if err != nil {
		return "", err
	}

	return url, nil
}
