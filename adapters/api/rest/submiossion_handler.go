package rest

import (
	"net/http"
	"qlass-be/dtos"
	"qlass-be/middleware"
	"qlass-be/usecases"
	"qlass-be/utils"

	"github.com/gin-gonic/gin"
)

type SubmissionHandler struct {
	SubmissionUseCase usecases.SubmissionUseCase
}

func NewSubmissionHandler(submissionUseCase usecases.SubmissionUseCase) *SubmissionHandler {
	return &SubmissionHandler{
		SubmissionUseCase: submissionUseCase,
	}
}

func (h *SubmissionHandler) CreateSubmission(c *gin.Context) {
	val, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, dtos.GlobalErrorResponse{Error: "UNAUTHORIZED", Message: "User context not found"})
		return
	}

	claims, ok := val.(*middleware.JWTCustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, dtos.GlobalErrorResponse{Error: "UNAUTHORIZED", Message: "Invalid user context"})
		return
	}

	var req dtos.CreateSubmissionDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.GlobalErrorResponse{Error: "BAD_REQUEST", Message: err.Error()})
		return
	}

	err := h.SubmissionUseCase.CreateSubmission(req, claims.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{Error: "INTERNAL_SERVER_ERROR", Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Submission created successfully"})
}

func (h *SubmissionHandler) GetSubmission(c *gin.Context) {
	submissionId := c.Param("id")

	val, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, dtos.GlobalErrorResponse{Error: "UNAUTHORIZED", Message: "User context not found"})
		return
	}

	_, ok := val.(*middleware.JWTCustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, dtos.GlobalErrorResponse{Error: "UNAUTHORIZED", Message: "Invalid user context"})
		return
	}

	submission, err := h.SubmissionUseCase.GetSubmissionByID(utils.StringToUint(submissionId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{Error: "INTERNAL_SERVER_ERROR", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, submission)
}

func (h *SubmissionHandler) GetSubmissonByMaterialIDAndStudentID(c *gin.Context) {
	class_material_id := c.Param("class_material_id")

	val, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, dtos.GlobalErrorResponse{Error: "UNAUTHORIZED", Message: "User context not found"})
		return
	}

	claims, ok := val.(*middleware.JWTCustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, dtos.GlobalErrorResponse{Error: "UNAUTHORIZED", Message: "Invalid user context"})
		return
	}

	submission, err := h.SubmissionUseCase.GetSubmissonByMaterialIDAndStudentID(utils.StringToUint(class_material_id), claims.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": submission})
}
