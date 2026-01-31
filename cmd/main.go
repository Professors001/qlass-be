package main

import (
	"context"
	"fmt"
	"log"

	"qlass-be/config"
	"qlass-be/domain/entities"
	"qlass-be/infrastructure/cache"
	"qlass-be/infrastructure/middleware"
	"qlass-be/infrastructure/storage"
	"qlass-be/router"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := config.NewPostgresDB(cfg)
	redisClient := config.NewRedisClient(cfg)
	jwtService := middleware.NewJWTService(cfg)
	minioClient := config.NewMinioClient(cfg)

	// Verify Redis connection
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v. Check REDIS_PASSWORD in .env", err)
	}

	cacheService := cache.NewCacheService(redisClient)
	cacheHelper := cache.NewCacheHelper(cacheService)

	// Verify MinIO Bucket
	exists, err := minioClient.BucketExists(context.Background(), cfg.MinioBucketName)
	if err != nil {
		log.Printf("⚠️  Error checking MinIO bucket: %v", err)
	} else if !exists {
		log.Printf("⚠️  MinIO bucket '%s' does not exist", cfg.MinioBucketName)
	} else {
		log.Printf("✅ MinIO bucket '%s' is ready", cfg.MinioBucketName)
	}

	storageService := storage.NewMinioStorageService(minioClient)

	// Migration
	if err := db.AutoMigrate(
		&entities.User{}, &entities.Class{}, &entities.ClassEnrollment{},
		&entities.CourseMaterial{}, &entities.Attachment{},
		&entities.Quiz{}, &entities.QuizQuestion{}, &entities.QuizOption{},
		&entities.QuizGameLog{}, &entities.Submission{}); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	// Init Gin Framework
	r := gin.Default()

	// Init Routers
	router.SetUpRouters(r, cfg, db, cacheHelper, jwtService, storageService) //, attachmentRepo

	serverAddr := fmt.Sprintf("%s", cfg.AppPort)
	log.Printf("🚀 Server running on port %s", serverAddr)
	r.Run(serverAddr)
}
