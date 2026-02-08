package rest

import (
	"net/http"
	"qlass-be/dtos"
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
	var req dtos.CreateQuizDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.GlobalErrorResponse{Error: "BAD_REQUEST", Message: err.Error()})
		return
	}

	quizID, err := h.UseCase.CreateQuiz(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{Error: "INTERNAL_ERROR", Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Quiz created successfully",
		"id":      quizID,
	})
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
