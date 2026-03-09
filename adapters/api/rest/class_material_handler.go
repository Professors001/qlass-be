package rest

import (
	"net/http"
	"qlass-be/dtos"
	"qlass-be/middleware"
	"qlass-be/usecases"
	"qlass-be/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MaterialHandler struct {
	materialUseCase usecases.ClassMaterialUseCase
}

func NewMaterialHandler(materialUseCase usecases.ClassMaterialUseCase) *MaterialHandler {
	return &MaterialHandler{
		materialUseCase: materialUseCase,
	}
}

func (h *MaterialHandler) CreateMaterial(c *gin.Context) {
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

	var req dtos.CreateClassMaterialDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.GlobalErrorResponse{Error: "BAD_REQUEST", Message: err.Error()})
		return
	}

	err := h.materialUseCase.CreateClassMaterial(&req, claims.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{Error: "INTERNAL_SERVER_ERROR", Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Material created successfully"})
}

func (h *MaterialHandler) GetMaterialByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.GlobalErrorResponse{Error: "BAD_REQUEST", Message: "Invalid material ID"})
		return
	}

	res, err := h.materialUseCase.GetMaterialByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{Error: "INTERNAL_ERROR", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *MaterialHandler) GetMaterialsByClassID(c *gin.Context) {
	id := c.Param("class_id")
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

	res, err := h.materialUseCase.GetMaterialsByClassID(utils.StringToUint(id), claims.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{Error: "INTERNAL_ERROR", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

// In adapters/api/rest/class_material_handler.go

func (h *MaterialHandler) CreateQuizMaterial(c *gin.Context) {
	var dto dtos.CreateQuizClassMaterialDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.materialUseCase.CreateQuizMaterial(dto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Quiz material created successfully"})
}

func (h *MaterialHandler) UpdatePostMaterial(c *gin.Context) {
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

	var dto dtos.UpdatePostClassMaterialDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, dtos.GlobalErrorResponse{Error: "BAD_REQUEST", Message: err.Error()})
		return
	}

	if err := h.materialUseCase.UpdatePostClassMaterial(&dto, claims.UserId); err != nil {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{Error: "INTERNAL_SERVER_ERROR", Message: err.Error()})
		return
	}

	res, err := h.materialUseCase.GetMaterialByID(dto.ClassMaterialID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{Error: "INTERNAL_ERROR", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *MaterialHandler) UpdateAssignmentMaterial(c *gin.Context) {
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

	var dto dtos.UpdateAssignmentClassMaterialDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, dtos.GlobalErrorResponse{Error: "BAD_REQUEST", Message: err.Error()})
		return
	}

	if err := h.materialUseCase.UpdateAssignmentClassMaterial(&dto, claims.UserId); err != nil {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{Error: "INTERNAL_SERVER_ERROR", Message: err.Error()})
		return
	}

	res, err := h.materialUseCase.GetMaterialByID(dto.ClassMaterialID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{Error: "INTERNAL_ERROR", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *MaterialHandler) UpdateQuizMaterial(c *gin.Context) {
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

	var dto dtos.UpdateQuizClassMaterialDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, dtos.GlobalErrorResponse{Error: "BAD_REQUEST", Message: err.Error()})
		return
	}

	if err := h.materialUseCase.UpdateQuizClassMaterial(&dto, claims.UserId); err != nil {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{Error: "INTERNAL_SERVER_ERROR", Message: err.Error()})
		return
	}
}
