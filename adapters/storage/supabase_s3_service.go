package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type supabaseS3StorageService struct {
	client  *s3.Client
	presign *s3.PresignClient
}

func NewSupabaseS3StorageService(endpoint, region, accessKey, secretKey string, useSSL bool) (StorageService, error) {
	baseEndpoint := strings.TrimSpace(endpoint)
	if baseEndpoint == "" {
		return nil, fmt.Errorf("MINIO_ENDPOINT is required")
	}

	if !strings.HasPrefix(baseEndpoint, "http://") && !strings.HasPrefix(baseEndpoint, "https://") {
		if useSSL {
			baseEndpoint = "https://" + baseEndpoint
		} else {
			baseEndpoint = "http://" + baseEndpoint
		}
	}

	if strings.Contains(baseEndpoint, "supabase.co") && !strings.Contains(baseEndpoint, "/storage/v1/s3") {
		baseEndpoint = strings.TrimRight(baseEndpoint, "/") + "/storage/v1/s3"
	}

	if strings.TrimSpace(region) == "" {
		region = "us-east-1"
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(baseEndpoint)
	})

	return &supabaseS3StorageService{
		client:  client,
		presign: s3.NewPresignClient(client),
	}, nil
}

func (s *supabaseS3StorageService) Upload(ctx context.Context, file *multipart.FileHeader, bucketName string, objectName string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectName),
		Body:        src,
		ContentType: aws.String(file.Header.Get("Content-Type")),
	})
	return err
}

func (s *supabaseS3StorageService) GetPresignedURL(ctx context.Context, bucketName string, objectName string, expiry time.Duration) (string, error) {
	out, err := s.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiry
	})
	if err != nil {
		return "", err
	}
	return out.URL, nil
}
