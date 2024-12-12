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
	"github.com/codevault-llc/minerva/pkg/logger"
	"go.uber.org/zap"
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

// ensureBucketExists checks if a bucket exists, creates it if not, and applies policies if necessary.
func ensureBucketExists(bucketName string) error {
	_, err := AWS.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		var noSuchBucket *types.NotFound
		if ok := errors.As(err, &noSuchBucket); ok {
			_, err = AWS.CreateBucket(context.TODO(), &s3.CreateBucketInput{
				Bucket: aws.String(bucketName),
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if bucketName == "content-bucket" {
		if err := setPublicBucketPolicy(bucketName); err != nil {
			return err
		}
	}

	return nil
}

// setPublicBucketPolicy applies a public access policy to the specified bucket.
func setPublicBucketPolicy(bucketName string) error {
	desiredPolicy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": "*",
				"Action": "s3:GetObject",
				"Resource": "arn:aws:s3:::` + bucketName + `/*"
			}
		]
	}`

	// Check the current policy
	currentPolicy, err := AWS.GetBucketPolicy(context.TODO(), &s3.GetBucketPolicyInput{
		Bucket: aws.String(bucketName),
	})
	if err == nil {
		// Compare existing policy with desired policy
		if *currentPolicy.Policy == desiredPolicy {
			return nil
		}
	} else {
		logger.Log.Info("Failed to get bucket policy for", zap.String("bucket", bucketName), zap.Error(err))
	}

	// Apply the new policy
	_, err = AWS.PutBucketPolicy(context.TODO(), &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucketName),
		Policy: aws.String(desiredPolicy),
	})
	if err != nil {
		return err
	}

	return nil
}
