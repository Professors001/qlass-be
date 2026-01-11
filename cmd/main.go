package main

import (
	"fmt"
	"log"

	"qlass-be/config"
	"qlass-be/domain/entities"
	"qlass-be/router"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := config.NewPostgresDB(cfg)

	// Migration
	db.AutoMigrate(
		&entities.User{}, &entities.Class{}, &entities.ClassEnrollment{},
		&entities.CourseMaterial{}, &entities.Attachment{},
		&entities.Quiz{}, &entities.QuizQuestion{}, &entities.QuizOption{},
		&entities.QuizGameLog{}, &entities.Submission{})

	// Init Gin Framework
	r := gin.Default()

	// Init Routers
	router.SetUpRouters(r, db)

	serverAddr := fmt.Sprintf("%s", cfg.AppPort)
	log.Printf("🚀 Server running on port %s", serverAddr)
	r.Run(serverAddr)
}
