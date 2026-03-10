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

	c.JSON(http.StatusOK, gin.H{"data": submission})
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

func (h *SubmissionHandler) GetStudentScores(c *gin.Context) {
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

	if claims.Role != "teacher" {
		c.JSON(http.StatusForbidden, dtos.GlobalErrorResponse{Error: "FORBIDDEN", Message: "Only teacher can view student score list"})
		return
	}

	var req dtos.TeacherGetStudentScoreListsByUserAndClassIDRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.GlobalErrorResponse{Error: "BAD_REQUEST", Message: err.Error()})
		return
	}

	studentScores, err := h.SubmissionUseCase.TeacherGetStudentScoreListsByUserAndClassID(req, claims.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{Error: "INTERNAL_SERVER_ERROR", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": studentScores})
}

func (h *SubmissionHandler) GetStudentSubmissionsByClass(c *gin.Context) {
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

	if claims.Role != "teacher" {
		c.JSON(http.StatusForbidden, dtos.GlobalErrorResponse{Error: "FORBIDDEN", Message: "Only teacher can view submissions by class"})
		return
	}

	req := dtos.GetStudentSubmissionsByClassRequestDto{ClassID: utils.StringToUint(c.Param("class_id"))}
	if req.ClassID == 0 {
		c.JSON(http.StatusBadRequest, dtos.GlobalErrorResponse{Error: "BAD_REQUEST", Message: "Invalid class ID"})
		return
	}

	submissions, err := h.SubmissionUseCase.GetStudentSubmissionsByClass(req, claims.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{Error: "INTERNAL_SERVER_ERROR", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": submissions})
}

func (h *SubmissionHandler) StudentSaveSubmission(c *gin.Context) {
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

	var req dtos.StudentSaveSubmissionDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.SubmissionUseCase.StudentSaveSubmission(req, claims.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	submission, err := h.SubmissionUseCase.GetSubmissionByID(req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{Error: "INTERNAL_SERVER_ERROR", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": submission})
}

func (h *SubmissionHandler) TeacherSaveSubmission(c *gin.Context) {
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

	var req dtos.TeacherSaveSubmissionDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.SubmissionUseCase.TeacherSaveSubmission(req, claims.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	submission, err := h.SubmissionUseCase.GetSubmissionByID(req.SubmissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{Error: "INTERNAL_SERVER_ERROR", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": submission})
}

func (h *SubmissionHandler) GetSubmissionsByMaterialID(c *gin.Context) {
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

	submissions, err := h.SubmissionUseCase.GetSubmissionsByMaterialID(utils.StringToUint(class_material_id), claims.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{Error: "INTERNAL_SERVER_ERROR", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": submissions})
}
