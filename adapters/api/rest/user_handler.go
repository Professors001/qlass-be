package rest

import (
	"net/http"
	"qlass-be/adapters/api/rest/utils"
	"qlass-be/dtos"
	"qlass-be/middleware"
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

func (h *UserHandler) CreateTeacher(c *gin.Context) {
	// Check Admin Role
	val, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, dtos.GlobalErrorResponse{Error: "UNAUTHORIZED", Message: "User context not found"})
		return
	}

	claims, ok := val.(*middleware.JWTCustomClaims)
	if !ok || claims.Role != "admin" {
		c.JSON(http.StatusForbidden, dtos.GlobalErrorResponse{Error: "FORBIDDEN", Message: "Only admin can perform this action"})
		return
	}

	// Bind DTO
	var req dtos.CreateTeacherRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.GlobalErrorResponse{Error: "BAD_REQUEST", Message: err.Error()})
		return
	}

	// Call UseCase
	res, err := h.UseCase.CreateTeacher(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.GlobalErrorResponse{Error: "CREATION_FAILED", Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	val, exists := c.Get("currentUser")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "No user found in context")
		return
	}

	claims, ok := val.(*middleware.JWTCustomClaims)
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid user context")
		return
	}

	var req dtos.UpdateUserRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	res, err := h.UseCase.UpdateUser(&req, claims.UserId)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	val, exists := c.Get("currentUser")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "No user found in context")
		return
	}

	claims, ok := val.(*middleware.JWTCustomClaims)
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid user context")
		return
	}

	var req dtos.ChangePasswordRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	res, err := h.UseCase.ChangePassword(&req, claims.UserId)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) ForgetPasswordStep1(c *gin.Context) {
	var req dtos.ForgetPasswordStep1RequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Call UseCase
	res, err := h.UseCase.ForgetPasswordStep1(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) ForgetPasswordStep2(c *gin.Context) {
	var req dtos.ForgetPasswordStep2RequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Call UseCase
	res, err := h.UseCase.ForgetPasswordStep2(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) AdminUpdateuser(c *gin.Context) {
	val, exists := c.Get("currentUser")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "No user found in context")
		return
	}

	claims, ok := val.(*middleware.JWTCustomClaims)
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid user context")
		return
	}

	if claims.Role != "admin" {
		utils.SendError(c, http.StatusForbidden, "FORBIDDEN", "Only admin can perform this action")
		return
	}

	var req dtos.AdminUpdateUserRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	err := h.UseCase.AdminUpdateUser(&req)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
