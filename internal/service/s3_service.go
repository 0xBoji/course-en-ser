package service

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

// S3Service handles S3 operations
type S3Service struct {
	client     *s3.S3
	bucketName string
	baseURL    string
	folder     string
}

// NewS3Service creates a new S3 service
func NewS3Service() *S3Service {
	// Get configuration from environment
	region := os.Getenv("S3_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("S3_BUCKET_NAME")
	baseURL := os.Getenv("S3_BASE_URL")
	folder := os.Getenv("S3_COURSE_IMAGES_FOLDER")

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKey,
			secretKey,
			"", // token
		),
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to create AWS session: %v", err))
	}

	return &S3Service{
		client:     s3.New(sess),
		bucketName: bucketName,
		baseURL:    baseURL,
		folder:     folder,
	}
}

// UploadCourseImage uploads a course image to S3
func (s *S3Service) UploadCourseImage(file *multipart.FileHeader) (string, error) {
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Validate file type
	if !isValidImageType(file.Filename) {
		return "", fmt.Errorf("invalid file type. Only JPG, JPEG, PNG, GIF, and WebP are allowed")
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return "", fmt.Errorf("file size too large. Maximum size is 5MB")
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	key := fmt.Sprintf("%s/%s", s.folder, filename)

	// Determine content type
	contentType := getContentType(ext)

	// Upload to S3
	_, err = s.client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        src,
		ContentType: aws.String(contentType),
		ACL:         aws.String("public-read"), // Make the file publicly accessible
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %v", err)
	}

	// Return the public URL
	imageURL := fmt.Sprintf("%s/%s", s.baseURL, key)
	return imageURL, nil
}

// DeleteCourseImage deletes a course image from S3
func (s *S3Service) DeleteCourseImage(imageURL string) error {
	// Extract key from URL
	key := strings.TrimPrefix(imageURL, s.baseURL+"/")

	// Delete from S3
	_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %v", err)
	}

	return nil
}

// isValidImageType checks if the file extension is a valid image type
func isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	return validTypes[ext]
}

// getContentType returns the appropriate content type for the file extension
func getContentType(ext string) string {
	ext = strings.ToLower(ext)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
