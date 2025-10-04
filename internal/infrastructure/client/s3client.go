package client

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"go.uber.org/zap"
	"io"
	"os"
	"v1/internal/infrastructure/config"
	"v1/internal/infrastructure/logger"

	cfgS3 "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// TODO Logger
type S3Client struct {
	client *s3.Client
	l      *zap.SugaredLogger
	bucket string
}

func NewS3Client(ctx context.Context, cfg config.Config) (*S3Client, error) {
	l := logger.GetLogger()
	configS3, err := cfgS3.LoadDefaultConfig(ctx, cfgS3.WithRegion(cfg.Region), cfgS3.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")))
	if err != nil {
		l.Errorf("failed to load AWS config: %v", err)
		return nil, fmt.Errorf("failed to load AWS config: %v", err)
	}

	client := s3.NewFromConfig(configS3)

	l.Infof("S3 client initialized: bucket=%s, region=%s", cfg.Bucket, cfg.Region)
	return &S3Client{
		client: client,
		l:      logger.GetLogger(),
		bucket: cfg.Bucket,
	}, nil
}

func (s *S3Client) UploadFile(ctx context.Context, key, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		s.l.Errorf("failed to open file %s: %v", filepath, err)
		return fmt.Errorf("failed to open file %s: %v", filepath, err)
	}
	defer file.Close()

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		s.l.Errorf("failed to upload file to S3: %v", err)
		return fmt.Errorf("failed to upload file to S3: %v", err)
	}

	s.l.Infof("S3 upload succeeded: bucket=%s, key=%s, local_path=%s", s.bucket, key, filepath)
	return nil
}

func (s *S3Client) DownloadFile(ctx context.Context, key, destination string) error {
	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		s.l.Errorf("failed to download file from S3: %v", err)
		return fmt.Errorf("failed to donload file from S3: %v", err)
	}
	defer resp.Body.Close()

	file, err := os.Create(destination)
	if err != nil {
		s.l.Errorf("failed to create file %s: %v", destination, err)
		return fmt.Errorf("failed to create file %s: %v", destination, err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		s.l.Errorf("failed to write file %s: %v", destination, err)
		return fmt.Errorf("failed to write file %s: %v", destination, err)
	}
	s.l.Infof("S3 download succeeded: bucket=%s, key=%s, destination=%s", s.bucket, key, destination)
	return nil
}

func (s *S3Client) ListObject(ctx context.Context) ([]string, error) {
	resp, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		s.l.Errorf("failed to list object: %v", err)
		return nil, fmt.Errorf("failed to list object: %v", err)
	}

	var keys []string
	for _, obj := range resp.Contents {
		keys = append(keys, *obj.Key)
	}

	s.l.Infof("S3 list objects succeeded: bucket=%s, object_count=%d", s.bucket, len(keys))
	return keys, nil
}
