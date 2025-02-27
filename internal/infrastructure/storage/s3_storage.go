package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"cams.dev/video_upload_backend/internal/domain/repository"
)

// S3Storage implements repository.StorageRepository using AWS S3 or Minio
type S3Storage struct {
	client     *s3.S3
	bucketName string
	region     string
	endpoint   string
}

// NewS3Storage creates a new S3 storage client
func NewS3Storage(
	accessKey, secretKey, region, bucketName, endpoint string,
	useSSL bool,
) repository.StorageRepository {
	// Configure AWS session
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		DisableSSL:       aws.Bool(!useSSL),
		S3ForcePathStyle: aws.Bool(true), // Required for Minio
	}

	sess, err := session.NewSession(s3Config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create S3 session: %v", err))
	}

	return &S3Storage{
		client:     s3.New(sess),
		bucketName: bucketName,
		region:     region,
		endpoint:   endpoint,
	}
}

// UploadFile uploads a file to S3/Minio
func (s *S3Storage) UploadFile(
	ctx context.Context,
	fileName string,
	data []byte,
	contentType string,
) (string, error) {
	// Create input for S3 PutObject
	input := &s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(fileName),
		Body:          bytes.NewReader(data),
		ContentLength: aws.Int64(int64(len(data))),
		ContentType:   aws.String(contentType),
	}

	// Upload the file
	_, err := s.client.PutObjectWithContext(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Build the URL
	var url string
	if s.endpoint != "" {
		// For Minio or custom S3 endpoint
		url = fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucketName, fileName)
	} else {
		// For AWS S3
		url = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.region, fileName)
	}

	return url, nil
}

// GetFile retrieves a file from S3/Minio
func (s *S3Storage) GetFile(ctx context.Context, fileName string) ([]byte, error) {
	// Create input for S3 GetObject
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
	}

	// Get the object
	result, err := s.client.GetObjectWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}
	defer result.Body.Close()

	// Read the data
	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file data: %w", err)
	}

	return data, nil
}

// GeneratePresignedURL generates a presigned URL for a file
func (s *S3Storage) GeneratePresignedURL(
	ctx context.Context,
	fileName string,
	expiry time.Duration,
) (string, error) {
	// Create a request for the presigned URL
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
	})

	// Generate the presigned URL
	url, err := req.Presign(expiry)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url, nil
}
