package rest

import (
	"net/http"
	"qlass-be/adapters/api/rest/utils"
	"qlass-be/dtos"
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

func (h *UserHandler) RegisterFirstStep(c *gin.Context) {
	// 1. Bind DTO
	var req dtos.RegisterRequestStepOneDto
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	// 2. Call UseCase
	res, err := h.UseCase.RegisterFirstStep(c.Request.Context(), &req)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "DUPLICATE_USER", err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) RegisterSecondStep(c *gin.Context) {
	// 1. Bind DTO
	var req dtos.RegisterRequestStepTwoDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Call UseCase
	res, err := h.UseCase.RegisterSecondStep(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// func (h *UserHandler) GetUser(c *gin.Context) {
// 	uuid := c.Param("uuid")

// 	user, err := h.UseCase.GetUserByUID(uuid) // Assumes you updated UseCase to use UUID
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

// 	// Clean conversion Domain -> DTO

// 	response := transform.ToUserResponse(user)
// 	c.JSON(http.StatusOK, response)
// }
