package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"qlass-be/adapters/cache"
	"qlass-be/adapters/storage"
	"qlass-be/config"
	"qlass-be/domain/entities"
	"qlass-be/middleware"
	"qlass-be/router"

	"github.com/gin-contrib/cors"
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
		&entities.ClassMaterial{}, &entities.Attachment{},
		&entities.Quiz{}, &entities.QuizQuestion{}, &entities.QuizOption{},
		&entities.QuizGameLog{}, &entities.Submission{}); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	// Init Gin Framework
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // อนุญาต Next.js ของคุณ
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // สำคัญมากถ้าจะใช้ Cookies/Auth.js
		MaxAge:           12 * time.Hour,
	}))

	// Init Routers
	router.SetUpRouters(r, cfg, db, cacheHelper, jwtService, storageService) //, attachmentRepo

	serverAddr := fmt.Sprintf("%s", cfg.AppPort)
	log.Printf("🚀 Server running on port %s", serverAddr)
	r.Run(serverAddr)
}
