package object

import (
	"fmt"

	"github.com/muzcategui1106/kitchen-wizard/pkg/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	// "github.com/aws/aws-sdk-go/aws/credentials"
)

// bucket constants
var (
	kitchenWizardBucket string = "kitchen-wizard"
)

// CreateClient creates connectivity to a minio tenant for object store
// TODO configure crendentials
func CreateClient(endpoint, accessKey, secretKey string) (*s3.S3, error) {
	// Set up the S3 client with generic credentials and endpoint
	cfg := aws.NewConfig().WithEndpoint(endpoint).WithRegion("us-east-1").WithLogLevel(aws.LogDebugWithHTTPBody)
	sess, err := session.NewSession(&aws.Config{
		// Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(true),
	})
	// sess, err := session.NewSessionWithOptions(cfg)
	if err != nil {
		logger.Log.Sugar().Errorf("error while connecting to s3 due to %s", err)
		return nil, err
	}

	svc := s3.New(sess, cfg)
	return svc, nil
}

// CreateNeccesaryBuckets creates necessary buckets for
func CreateNeccesaryBuckets(s3Client *s3.S3) error {
	buckets, err := s3Client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		logger.Log.Sugar().Errorf("could not list buckets due to %s", err)
	}
	fmt.Println(buckets)

	for _, b := range buckets.Buckets {
		if *b.Name == kitchenWizardBucket {
			return nil
		}
	}

	logger.Log.Sugar().Infof("%s bucket does not exist. will create it", kitchenWizardBucket)
	_, err = s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(kitchenWizardBucket),
	})

	if err != nil {
		logger.Log.Sugar().Errorf("could not create bucket due to %s", err)
		return err
	}

	return nil
}
