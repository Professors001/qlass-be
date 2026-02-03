package rest

import (
	"net/http"
	"qlass-be/dtos"
	"qlass-be/middleware"
	"qlass-be/usecases"

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
