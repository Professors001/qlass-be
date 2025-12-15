package rest

import (
	"net/http"
	"qlass-be/domain/dto/request"
	"qlass-be/domain/dto/response"
	"qlass-be/domain/entities"
	"qlass-be/usecases"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UseCase usecases.UserUseCase
}

func NewUserHandler(userUseCase usecases.UserUseCase) *UserHandler {
	return &UserHandler{
		UseCase: userUseCase,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	// 1. Bind DTO
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Map DTO -> Domain
	user := entities.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	// 3. Call Logic
	if err := h.UseCase.Register(&user, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Respond with DTO (Hide sensitive data)
	response := response.ToUserResponse(&user)
	c.JSON(http.StatusCreated, response)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	uuid := c.Param("uuid")

	user, err := h.UseCase.GetUserByUUID(uuid) // Assumes you updated UseCase to use UUID
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Clean conversion Domain -> DTO
	
	response := response.ToUserResponse(user)
	c.JSON(http.StatusOK, response)
}