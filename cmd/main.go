package main

import (
	"fmt"
	"log"

	"qlass-be/config"
	"qlass-be/domain/entities"
	"qlass-be/rounter"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := config.NewPostgresDB(cfg)

	// Migration
	db.AutoMigrate(&entities.User{})

	// Init Gin Framework
	r := gin.Default()

	// Init Routers
	rounter.SetUpRouters(r, db)

	serverAddr := fmt.Sprintf("%s", cfg.AppPort)
	log.Printf("🚀 Server running on port %s", serverAddr)
	r.Run(serverAddr)
}
