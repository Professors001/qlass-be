package router

import (
	"net/http"
	"qlass-be/adapters/api"
	"qlass-be/adapters/api/rest"
	"qlass-be/adapters/databases"
	"qlass-be/adapters/redis"
	"qlass-be/infrastructure/cache"
	"qlass-be/infrastructure/middleware"
	"qlass-be/usecases"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUpRouters(r *gin.Engine, db *gorm.DB, cache *cache.CacheHelper, jwtService middleware.JwtService) {

	// Seed Admin User if not exists
	databases.SeedAdminUser(db)

	userRepo := databases.NewPostgresUserRepository(db)
	userCacheRepo := redis.NewUserRedisRepository(cache)
	userUseCase := usecases.NewUserUseCase(userRepo, userCacheRepo, jwtService)
	userHandler := rest.NewUserHandler(userUseCase)

	classRepo := databases.NewPostgresClassRepository(db)
	classUseCase := usecases.NewClassUseCase(classRepo)
	classHandler := rest.NewClassHandler(classUseCase)

	handler := api.ProvideHandler(userHandler, classHandler)

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
	classRouter.POST("/", middleware.AuthorizeJWT(jwtService), handler.ClassHandler.CreateClass)
	classRouter.GET("/", middleware.AuthorizeJWT(jwtService), handler.ClassHandler.GetAllMyClasses)
	classRouter.GET("/invite/:code", handler.ClassHandler.GetClassDetailsByInviteCode)
}
