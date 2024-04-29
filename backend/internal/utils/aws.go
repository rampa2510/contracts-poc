package utils

import (
	"context"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AwsClient struct {
	presignClient *s3.PresignClient
}

func NewAwsClient(awsRegion string) *AwsClient {
	slog.Info("Initialising Aws Client")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		slog.Error("Error", err)
		panic(err)
	}

	s3Client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(s3Client)

	slog.Info("Initialised aws client")

	return &AwsClient{presignClient: presignClient}
}

func (client *AwsClient) GetPresignedUrl(bucketName, key string) string {
	presignedPutResponse, err := client.presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(3600 * int64(time.Second))
	})
	if err != nil {
		return ""
	}

	slog.Info("Generated presigned URL", "url", presignedPutResponse.URL)
	return presignedPutResponse.URL
}
