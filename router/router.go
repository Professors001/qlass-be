package router

import (
	"net/http"
	"qlass-be/adapters/api"
	"qlass-be/adapters/api/rest"
	"qlass-be/adapters/cache"
	"qlass-be/adapters/databases"
	"qlass-be/adapters/storage"
	"qlass-be/config"
	"qlass-be/middleware"
	"qlass-be/usecases"
	_ "qlass-be/usecases"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUpRouters(r *gin.Engine, cfg *config.Config, db *gorm.DB, cacheService *cache.CacheHelper, jwtService middleware.JwtService, storageService storage.StorageService) {

	// Seed Users (Admin, Teacher, Student) if not exists
	databases.SeedUsers(db)

	userRepo := databases.NewPostgresUserRepository(db)
	userCacheRepo := cache.NewUserRedisRepository(cacheService)
	userUseCase := usecases.NewUserUseCase(userRepo, userCacheRepo, jwtService)
	userHandler := rest.NewUserHandler(userUseCase)

	classRepo := databases.NewPostgresClassRepository(db)
	enrollRepo := databases.NewPostgresEnrollRepository(db)
	classUseCase := usecases.NewClassUseCase(classRepo, enrollRepo)
	classHandler := rest.NewClassHandler(classUseCase)

	attachmentRepo := databases.NewPostgresAttachmentRepository(db)
	attachmentUseCase := usecases.NewAttachmentUseCase(storageService, attachmentRepo, userRepo, cfg)
	attachmentHandler := rest.NewAttachmentHandler(attachmentUseCase)

	classMaterialRepo := databases.NewPostgresClassMaterialRepository(db)
	classMaterialUseCase := usecases.NewClassMaterialUseCase(classMaterialRepo, classRepo, attachmentRepo)
	classMaterialHandler := rest.NewMaterialHandler(classMaterialUseCase)

	handler := api.ProvideHandler(
		userHandler,
		classHandler,
		attachmentHandler,
		classMaterialHandler)
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Qlass BE is still running!"})
	})

	// Users
	userRouter := r.Group("/auth")
	userRouter.POST("/register-step-one", handler.UserHandler.RegisterFirstStep)
	userRouter.POST("/register-step-two", handler.UserHandler.RegisterSecondStep)
	userRouter.POST("/login", handler.UserHandler.Login)
	userRouter.GET("/me", middleware.AuthorizeJWT(jwtService), handler.UserHandler.Me)
	userRouter.POST("/create-teacher", middleware.AuthorizeJWT(jwtService), handler.UserHandler.CreateTeacher)

	// Classes
	classRouter := r.Group("/classes")
	classRouter.POST("", middleware.AuthorizeJWT(jwtService), handler.ClassHandler.CreateClass)
	classRouter.GET("", middleware.AuthorizeJWT(jwtService), handler.ClassHandler.GetAllMyClasses)
	classRouter.GET("/:id", middleware.AuthorizeJWT(jwtService), handler.ClassHandler.GetClassByID)
	classRouter.POST("/join", middleware.AuthorizeJWT(jwtService), handler.ClassHandler.EnrollStudent)
	classRouter.GET("/:id/students", middleware.AuthorizeJWT(jwtService), handler.ClassHandler.GetEnrolledStudents)
	classRouter.GET("/invite/:code", handler.ClassHandler.GetClassDetailsByInviteCode)

	// Attachments
	attachmentRouter := r.Group("/attachments")
	attachmentRouter.Use(middleware.AuthorizeJWT(jwtService))
	attachmentRouter.POST("", handler.AttachmentHandler.UploadAttachment)
	attachmentRouter.GET("/:attachmentID", handler.AttachmentHandler.GetAttachment)

	// Class Materials
	materialRouter := r.Group("/materials")
	materialRouter.Use(middleware.AuthorizeJWT(jwtService))
	materialRouter.POST("", handler.ClassMaterialHandler.CreateMaterial)
}
