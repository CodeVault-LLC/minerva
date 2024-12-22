package storage

import (
	"bytes"
	"context"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/codevault-llc/minerva/internal/database"
	"github.com/codevault-llc/minerva/pkg/utils"
	"github.com/google/uuid"
)

// DetermineStorageType decides whether content should be in hot or cold storage.
func DetermineStorageType(content string) string {
	if len(content) < 1024*100 {
		return "hot"
	}
	return "cold"
}

var supportedFileExtensions = []string{"jpg", "jpeg", "png", "pdf", "txt", "html", "js", "mjs", "mp3", "mp4", "css", "scss"}

// GetFileExtension returns the file extension of a given file name.
func GetFileExtension(fileName string) string {
	if fileName == "" {
		return "txt"
	}

	parts := strings.Split(fileName, ".")
	if len(parts) == 1 {
		return "txt"
	}

	if !utils.StringInSlice(parts[len(parts)-1], supportedFileExtensions) {
		return "txt"
	}

	return parts[len(parts)-1]
}

// SanitizeObjectKey removes unsupported characters from the object key.
func SanitizeObjectKey(key string) string {
	key = strings.ReplaceAll(key, " ", "_")
	re := regexp.MustCompile(`[^a-zA-Z0-9._/-]+`)
	key = re.ReplaceAllString(key, "")
	return key
}

func GetContentType(fileExtension string) string {
	switch fileExtension {
	case "jpg", "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "pdf":
		return "application/pdf"
	case "txt":
		return "text/plain"
	case "html":
		return "text/html"
	default:
		return "application/octet-stream"
	}
}

// GenerateObjectKey creates a unique object key with the file's original extension.
func GenerateObjectKey(originalFileName string) string {
	ext := GetFileExtension(originalFileName)
	id := uuid.New().String()
	return id + "." + ext
}

func UploadFile(bucketName string, objectKey string, fileContents []byte, contentType string, readable bool) error {
	input := &s3.PutObjectInput{
		Bucket:      &bucketName,
		Key:         &objectKey,
		Body:        bytes.NewReader(fileContents),
		ContentType: aws.String(contentType),
	}

	if readable {
		input.ACL = types.ObjectCannedACLPublicRead
	}

	_, err := database.AWS.PutObject(context.TODO(), input)
	return err
}

func DownloadFile(bucketName string, objectKey string) ([]byte, error) {
	resp, err := database.AWS.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DeleteFile(bucketName string, objectKey string) error {
	_, err := database.AWS.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})
	return err
}

func ListFiles(bucketName string) ([]string, error) {
	resp, err := database.AWS.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	})
	if err != nil {
		return nil, err
	}

	var keys []string
	for _, obj := range resp.Contents {
		keys = append(keys, *obj.Key)
	}

	return keys, nil
}

func GetEndpoint(bucketName string) string {
	return "http://localhost:9000/" + bucketName
}

func GetLocation(bucketName, objectKey string) string {
	return "http://localhost:9000/" + bucketName + "/" + objectKey
}
