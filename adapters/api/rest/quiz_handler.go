package rest

import (
	"net/http"
	"qlass-be/dtos"
	"qlass-be/middleware"
	"qlass-be/usecases"
	"qlass-be/utils"

	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	UseCase usecases.QuizUseCase
}

func NewQuizHandler(useCase usecases.QuizUseCase) *QuizHandler {
	return &QuizHandler{
		UseCase: useCase,
	}
}

func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	val, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, dtos.GlobalErrorResponse{Error: "UNAUTHORIZED", Message: "User context not found"})
		return
	}

	claims, ok := val.(*middleware.JWTCustomClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: "Failed to cast user context",
		})
		return
	}

	if claims.Role != "teacher" {
		c.JSON(http.StatusForbidden, dtos.GlobalErrorResponse{
			Error:   "FORBIDDEN",
			Message: "Only teachers can create classes",
		})
		return
	}

	var req dtos.SaveQuizDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.GlobalErrorResponse{Error: "BAD_REQUEST", Message: err.Error()})
		return
	}

	res, err := h.UseCase.CreateQuiz(req, claims.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{Error: "INTERNAL_ERROR", Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Quiz created successfully", "data": res})
}

func (h *QuizHandler) UpdateQuiz(c *gin.Context) {
	val, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, dtos.GlobalErrorResponse{Error: "UNAUTHORIZED", Message: "User context not found"})
		return
	}

	claims, ok := val.(*middleware.JWTCustomClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: "Failed to cast user context",
		})
		return
	}

	if claims.Role != "teacher" {
		c.JSON(http.StatusForbidden, dtos.GlobalErrorResponse{
			Error:   "FORBIDDEN",
			Message: "Only teachers can create classes",
		})
		return
	}

	var req dtos.SaveQuizDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.GlobalErrorResponse{Error: "BAD_REQUEST", Message: err.Error()})
		return
	}

	err := h.UseCase.UpdateQuiz(req, claims.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{Error: "INTERNAL_ERROR", Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Quiz updated successfully"})
}

func (h *QuizHandler) SaveQuizQuestion(c *gin.Context) {
	val, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, dtos.GlobalErrorResponse{Error: "UNAUTHORIZED", Message: "User context not found"})
		return
	}

	claims, ok := val.(*middleware.JWTCustomClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: "Failed to cast user context",
		})
		return
	}

	if claims.Role != "teacher" {
		c.JSON(http.StatusForbidden, dtos.GlobalErrorResponse{
			Error:   "FORBIDDEN",
			Message: "Only teachers can create classes",
		})
		return
	}

	quiz_id := utils.StringToUint(c.Param("id"))

	_, err := h.UseCase.GetQuizByID(quiz_id)
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.GlobalErrorResponse{Error: "NOT_FOUND", Message: err.Error()})
		return
	}

	var req dtos.SaveQuizQuestionDtoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.GlobalErrorResponse{Error: "BAD_REQUEST", Message: err.Error()})
		return
	}

	err = h.UseCase.SaveQuizQuestion(req, quiz_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{Error: "INTERNAL_ERROR", Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Quiz updated successfully"})
}

func (h *QuizHandler) GetQuiz(c *gin.Context) {
	id := utils.StringToUint(c.Param("id"))

	res, err := h.UseCase.GetQuizByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.GlobalErrorResponse{Error: "NOT_FOUND", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
