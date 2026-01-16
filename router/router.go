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

	userRepo := databases.NewPostgresUserRepository(db)
	userCacheRepo := redis.NewUserRedisRepository(cache)
	userUseCase := usecases.NewUserUseCase(userRepo, userCacheRepo, jwtService)
	userHandler := rest.NewUserHandler(userUseCase)

	handler := api.ProvideHandler(userHandler)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Qlass BE is still running!"})
	})

	// Users
	userRouter := r.Group("/users")
	userRouter.POST("/register-step-one", handler.UserHandler.RegisterFirstStep)
	userRouter.POST("/register-step-two", handler.UserHandler.RegisterSecondStep)
	userRouter.POST("/login", handler.UserHandler.Login)
	userRouter.GET("/me", middleware.AuthorizeJWT(jwtService), handler.UserHandler.Me)
}
