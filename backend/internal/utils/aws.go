package utils

import (
	"bytes"
	"context"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AwsClient struct {
	presignClient *s3.PresignClient
	s3Client      *s3.Client
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

	return &AwsClient{presignClient: presignClient, s3Client: s3Client}
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

func (client *AwsClient) S3Download(bucket, key string) (*s3.GetObjectOutput, error) {
	s3ObjectResp, err := client.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return s3ObjectResp, nil
}

func (client *AwsClient) S3Upload(bucket string, key string, data []byte) error {
	_, err := client.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return err
	}
	return nil
}
