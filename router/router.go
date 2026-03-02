package router

import (
	"net/http"
	"qlass-be/adapters/api"
	"qlass-be/adapters/api/rest"
	"qlass-be/adapters/api/websocket"
	"qlass-be/adapters/cache"
	"qlass-be/adapters/databases"
	"qlass-be/adapters/services"
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
	emailService := services.NewSMTPEmailService(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUser,
		cfg.SMTPPass,
	)
	attachmentRepo := databases.NewPostgresAttachmentRepository(db)

	userRepo := databases.NewPostgresUserRepository(db)
	userCacheRepo := cache.NewUserRedisRepository(cacheService)

	attachmentUseCase := usecases.NewAttachmentUseCase(storageService, attachmentRepo, userRepo, cfg)
	attachmentHandler := rest.NewAttachmentHandler(attachmentUseCase)

	userUseCase := usecases.NewUserUseCase(userRepo, userCacheRepo, jwtService, emailService, attachmentUseCase)
	userHandler := rest.NewUserHandler(userUseCase)

	classRepo := databases.NewPostgresClassRepository(db)
	enrollRepo := databases.NewPostgresEnrollRepository(db)
	classUseCase := usecases.NewClassUseCase(classRepo, enrollRepo)
	classHandler := rest.NewClassHandler(classUseCase)

	quizRepo := databases.NewPostgresQuizRepository(db)
	quizGameLogRepo := databases.NewPostgresQuizGameLogRepository(db)

	classMaterialRepo := databases.NewPostgresClassMaterialRepository(db)
	classMaterialUseCase := usecases.NewClassMaterialUseCase(classMaterialRepo, classRepo, attachmentRepo, attachmentUseCase, quizGameLogRepo, quizRepo)
	classMaterialHandler := rest.NewMaterialHandler(classMaterialUseCase)

	submissionRepo := databases.NewPostgresSubmissionRepository(db)
	submissionUseCase := usecases.NewSubmissionUseCase(submissionRepo, classMaterialRepo, attachmentRepo, attachmentUseCase)
	submissionHandler := rest.NewSubmissionHandler(submissionUseCase)

	quizQuestionRepo := databases.NewPostgresQuizQuestionRepository(db)
	quizOptionRepo := databases.NewPostgresQuizOptionRepository(db)
	quizStudentResponseRepo := databases.NewPostgresQuizStudentResponseRepository(db)
	quizUseCase := usecases.NewQuizUseCase(quizRepo, quizQuestionRepo, quizOptionRepo, attachmentRepo, attachmentUseCase)
	quizHandler := rest.NewQuizHandler(quizUseCase)

	gameRepo := cache.NewGameRedisRepository(cacheService)
	gameUseCase := usecases.NewGameUseCase(gameRepo, quizGameLogRepo, classMaterialRepo, classRepo, userRepo, submissionRepo, quizStudentResponseRepo)
	gameHandler := rest.NewGameHandler(gameUseCase)

	handler := api.ProvideHandler(
		userHandler,
		classHandler,
		attachmentHandler,
		classMaterialHandler,
		submissionHandler,
		quizHandler)
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
	userRouter.PUT("/update", middleware.AuthorizeJWT(jwtService), handler.UserHandler.UpdateUser)
	userRouter.PUT("/change-password", middleware.AuthorizeJWT(jwtService), handler.UserHandler.ChangePassword)
	userRouter.POST("/forgot-password-step-one", handler.UserHandler.ForgetPasswordStep1)
	userRouter.POST("/forgot-password-step-two", handler.UserHandler.ForgetPasswordStep2)
	userRouter.PUT("/admin-update", middleware.AuthorizeJWT(jwtService), handler.UserHandler.AdminUpdateuser)

	// Classes
	classRouter := r.Group("/classes")
	classRouter.POST("", middleware.AuthorizeJWT(jwtService), handler.ClassHandler.CreateClass)
	classRouter.GET("", middleware.AuthorizeJWT(jwtService), handler.ClassHandler.GetAllMyClasses)
	classRouter.GET("/:id", middleware.AuthorizeJWT(jwtService), handler.ClassHandler.GetClassByID)
	classRouter.POST("/join", middleware.AuthorizeJWT(jwtService), handler.ClassHandler.EnrollStudent)
	classRouter.GET("/:id/students", middleware.AuthorizeJWT(jwtService), handler.ClassHandler.GetEnrolledStudents)
	classRouter.GET("/invite/:code", handler.ClassHandler.GetClassDetailsByInviteCode)
	classRouter.PUT("/update", middleware.AuthorizeJWT(jwtService), handler.ClassHandler.UpdateClass)

	// Attachments
	attachmentRouter := r.Group("/attachments")
	attachmentRouter.Use(middleware.AuthorizeJWT(jwtService))
	attachmentRouter.POST("/upload", handler.AttachmentHandler.UploadAttachment)
	attachmentRouter.GET("/:attachmentID", handler.AttachmentHandler.GetAttachment)

	// Class Materials
	materialRouter := r.Group("/materials")
	materialRouter.Use(middleware.AuthorizeJWT(jwtService))
	materialRouter.POST("", handler.ClassMaterialHandler.CreateMaterial)
	materialRouter.GET("/:id", handler.ClassMaterialHandler.GetMaterialByID)
	materialRouter.POST("/quiz", handler.ClassMaterialHandler.CreateQuizMaterial)
	materialRouter.GET("/class/:class_id", handler.ClassMaterialHandler.GetMaterialsByClassID)

	// Submissions
	submissionRouter := r.Group("/submissions")
	submissionRouter.Use(middleware.AuthorizeJWT(jwtService))
	submissionRouter.POST("", handler.SubmissionHandler.CreateSubmission)
	submissionRouter.GET("/:id", handler.SubmissionHandler.GetSubmission)
	submissionRouter.GET("/material/:class_material_id", handler.SubmissionHandler.GetSubmissonByMaterialIDAndStudentID)

	// Quizzes
	quizRouter := r.Group("/quizzes")
	quizRouter.Use(middleware.AuthorizeJWT(jwtService))
	quizRouter.POST("", handler.QuizHandler.CreateQuiz)
	quizRouter.PUT("/:id", handler.QuizHandler.UpdateQuiz)
	quizRouter.POST("/:id/questions", handler.QuizHandler.SaveQuizQuestion)
	quizRouter.GET("/:id", handler.QuizHandler.GetQuiz)
	quizRouter.GET("/user", handler.QuizHandler.GetQuizzesByUserID)

	// Games
	gameRouter := r.Group("/games")
	gameRouter.Use(middleware.AuthorizeJWT(jwtService))
	gameRouter.POST("/start", gameHandler.StartGame)

	// Websocket

	manager := websocket.NewManager(gameUseCase)
	r.GET("/ws", func(c *gin.Context) {
		// 1. Get Token
		tokenString := c.Query("token")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}

		// 2. Get PIN (New)
		pin := c.Query("pin")
		if pin == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing game pin"})
			return
		}

		// 3. Validate Token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// 4. Pass PIN to ServeWS
		manager.ServeWS(c.Writer, c.Request, claims, pin)
	})
}
