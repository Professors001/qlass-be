package router

import (
	"net/http"
	"qlass-be/adapters/api"
	"qlass-be/adapters/api/rest"
	"qlass-be/adapters/databases"
	"qlass-be/usecases"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUpRouters(r *gin.Engine, db *gorm.DB) {

	userRepo := databases.NewPostgresUserRepository(db)
	userUseCase := usecases.NewUserUseCase(userRepo)
	userHandler := rest.NewUserHandler(userUseCase)

	handler := api.ProvideHandler(userHandler)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Qlass BE is still running!"})
	})

	// Users
	userRouter := r.Group("/users")
	userRouter.POST("/register", handler.UserHandler.Register)
	userRouter.GET("/:uuid", handler.UserHandler.GetUser)
}
