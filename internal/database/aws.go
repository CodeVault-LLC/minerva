package database

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
)

var AWS *s3.Client

type awsEndpointResolver struct {
}

func (*awsEndpointResolver) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (
	smithyendpoints.Endpoint, error,
) {
	return s3.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

// InitAWS initializes the AWS S3 client and ensures specified buckets exist.
func InitAWS() error {
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion("us-east-1"),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			"admin",    // Access Key ID (matches MINIO_ROOT_USER)
			"admin123", // Secret Access Key (matches MINIO_ROOT_PASSWORD)
			"",         // Session Token (not needed for MinIO)
		)),
	)

	if err != nil {
		log.Printf("Error loading AWS config: %v", err)
		return err
	}

	AWS = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // Necessary for MinIO compatibility
		o.BaseEndpoint = aws.String("http://localhost:9000")
		o.EndpointResolverV2 = &awsEndpointResolver{}
	})

	// Test the connection by listing buckets.
	if _, err := AWS.ListBuckets(context.TODO(), &s3.ListBucketsInput{}); err != nil {
		log.Printf("Failed to connect to S3: %v", err)
		return err
	}

	// Ensure necessary buckets exist.
	requiredBuckets := []string{"content-bucket", "logs-bucket"}
	for _, bucket := range requiredBuckets {
		if err := ensureBucketExists(bucket); err != nil {
			log.Printf("Error ensuring bucket %s exists: %v", bucket, err)
			return err
		}
	}

	return nil
}

// ensureBucketExists checks if a bucket exists and creates it if not.
func ensureBucketExists(bucketName string) error {
	_, err := AWS.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		var noSuchBucket *types.NotFound
		if ok := errors.As(err, &noSuchBucket); ok {
			// Create the bucket since it doesn't exist.
			_, err = AWS.CreateBucket(context.TODO(), &s3.CreateBucketInput{
				Bucket: aws.String(bucketName),
			})
			if err != nil {
				return err
			}
			log.Printf("Bucket %s created successfully", bucketName)
		} else {
			// If it's a different error, return it.
			return err
		}
	}

	return nil
}
