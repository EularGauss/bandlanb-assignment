package s3service

import (
	"fmt"
	"sync"
	"github.com/EularGauss/bandlab-assignment/internal/app"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3ServiceImpl struct {
	Bucket string
}

var s3Service *S3ServiceImpl
var mu sync.Mutex = sync.Mutex{}

func GetS3Service() *S3ServiceImpl {
	mu.Lock()
	config := DefaultS3Config()
	if s3Service == nil {
		s3Service = &S3ServiceImpl{
			Bucket: DefaultS3Config(),
		}
	}
	mu.Unlock()
	return s3Service
}

// GeneratePresignedURL generates a pre-signed URL for uploading an image to S3
func (s *S3ServiceImpl) GeneratePresignedURL(key string) (string, error) {
	// Check for allowed file extensions
	allowedExtensions := []string{".jpg", ".jpeg", ".bmp"}
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
	url, err := req.Presign(time.Duration(s3_config.PreSignedURLExpiration) * time.Minute)
	if err != nil {
		return "", err
	}

	return url, nil
}
