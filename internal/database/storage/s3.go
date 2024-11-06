package storage

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/codevault-llc/minerva/internal/database"
)

// DetermineStorageType decides whether content should be in hot or cold storage.
func DetermineStorageType(content string) string {
	if len(content) < 1024*100 {
		return "hot"
	}
	return "cold"
}

func UploadFile(bucketName string, objectKey string, fileContents []byte, readable bool) error {
	if readable {
		_, err := database.AWS.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: &bucketName,
			Key:    &objectKey,
			Body:   bytes.NewReader(fileContents),
			ACL:    types.ObjectCannedACLPublicRead,
		})

		return err
	} else {
		_, err := database.AWS.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: &bucketName,
			Key:    &objectKey,
			Body:   bytes.NewReader(fileContents),
		})

		return err
	}
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
