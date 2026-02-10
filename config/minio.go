package config

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioClient(cfg *Config) *minio.Client {
	// Initialize MinIO client object
	minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		log.Fatalf("❌ Failed to create MinIO client: %v", err)
	}

	// Verify connection by listing buckets (lightweight check)
	if _, err := minioClient.ListBuckets(context.Background()); err != nil {
		log.Printf("⚠️  Failed to connect to MinIO at %s: %v", cfg.MinioEndpoint, err)
	} else {
		log.Println("✅ Connected to MinIO successfully!")
	}

	return minioClient
}
