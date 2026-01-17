package rest

import (
	"net/http"
	"qlass-be/adapters/api/rest/error"
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

func (h *UserHandler) Login(c *gin.Context) {
	// 1. Bind DTO
	var req dtos.LoginRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	// 2. Call UseCase
	res, err := h.UseCase.Login(c.Request.Context(), &req)
	if err != nil {
		utils.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) Me(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "No user found in context")
		return
	}
	c.JSON(http.StatusOK, user)
}
