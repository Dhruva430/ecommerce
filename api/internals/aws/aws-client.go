package awsclient

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var client *s3.Client

// GetS3Client initializes the S3 client
func GetS3Client() *s3.Client {
	if client != nil {
		return client
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		log.Fatal("unable to load AWS config:", err)
	}

	client = s3.NewFromConfig(cfg)
	return client
}

// ListObjects lists objects in the bucket
func ListObjects(client *s3.Client) {
	bucket := os.Getenv("AWS_S3_BUCKET")
	if bucket == "" {
		log.Fatal("AWS_S3_BUCKET not set")
	}

	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Objects in bucket:")
	for _, object := range output.Contents {
		log.Printf("key=%s size=%d", aws.ToString(object.Key), *object.Size)
	}
}

func GeneratePresignedUploadURL(key string, contentType string, fileSize int64) (string, error) {
	client := GetS3Client()
	presignedClient := s3.NewPresignClient(client)
	bucket := os.Getenv("AWS_S3_BUCKET")
	if bucket == "" {
		return "", nil
	}

	params := &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(fileSize),
	}

	presignedReq, err := presignedClient.PresignPutObject(context.TODO(), params)
	if err != nil {
		return "", err
	}
	return presignedReq.URL, nil
}
