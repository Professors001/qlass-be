package config

import (
	"context"
	"log"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioClient(cfg *Config) *minio.Client {
	endpoint := strings.TrimSpace(cfg.MinioEndpoint)
	if endpoint == "" {
		log.Fatal("❌ MINIO_ENDPOINT is required")
	}

	// minio-go expects host[:port] only. If a full URL is provided, normalize it.
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		u, err := url.Parse(endpoint)
		if err != nil || u.Host == "" {
			log.Fatalf("❌ Invalid MINIO_ENDPOINT: %s", cfg.MinioEndpoint)
		}
		if u.Path != "" && u.Path != "/" {
			log.Printf("⚠️  MINIO_ENDPOINT path '%s' is ignored by minio-go. Using host only: %s", u.Path, u.Host)
		}
		endpoint = u.Host
	}

	// Initialize MinIO client object
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		log.Fatalf("❌ Failed to create MinIO client: %v", err)
	}

	// Supabase S3 is compatible for object operations, but MinIO bucket probing can be misleading.
	if strings.Contains(endpoint, "supabase.co") {
		log.Println("✅ Storage client initialized for Supabase S3")
		return minioClient
	}

	// Verify connection by listing buckets (lightweight check)
	if _, err := minioClient.ListBuckets(context.Background()); err != nil {
		log.Printf("⚠️  Failed to connect to MinIO at %s: %v", endpoint, err)
	} else {
		log.Println("✅ Connected to MinIO successfully!")
	}

	return minioClient
}
