package rest

import (
	"net/http"
	"qlass-be/dtos"
	"qlass-be/infrastructure/middleware" // ตรวจสอบว่าชี้ไปยัง package ที่มี JWTCustomClaims
	"qlass-be/usecases"

	"github.com/gin-gonic/gin"
)

type ClassHandler struct {
	UseCase usecases.ClassUseCase
}

func NewClassHandler(useCase usecases.ClassUseCase) *ClassHandler {
	return &ClassHandler{
		UseCase: useCase,
	}
}

func (h *ClassHandler) CreateClass(c *gin.Context) {
	// 1. Extraction: ดึง Claims จาก Context
	val, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, dtos.GlobalErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "User context not found",
		})
		return
	}

	// 2. Type Assertion: แปลง Interface เป็น JWTCustomClaims (ตัวล่าสุดที่เป็น uint)
	claims, ok := val.(*middleware.JWTCustomClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: "Failed to cast user context",
		})
		return
	}

	// 3. Parsing: รับข้อมูล JSON Request Body
	var req dtos.CreateClassRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.GlobalErrorResponse{
			Error:   "INVALID_INPUT",
			Message: err.Error(),
		})
		return
	}

	// 4. Execution: เรียกใช้งาน UseCase โดยส่ง UserId (uint) ไปได้ทันที
	// Senior Note: ไม่ต้องใช้ strconv.ParseUint แล้ว เพราะข้อมูลเป็น uint ตั้งแต่ต้นทาง
	res, err := h.UseCase.CreateClass(c.Request.Context(), &req, claims.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	// 5. Response
	c.JSON(http.StatusCreated, res)
}

func (h *ClassHandler) GetClassDetailsByInviteCode(c *gin.Context) {
	inviteCode := c.Param("code")
	if inviteCode == "" {
		c.JSON(http.StatusBadRequest, dtos.GlobalErrorResponse{
			Error:   "BAD_REQUEST",
			Message: "Invite code is required",
		})
		return
	}

	res, err := h.UseCase.GetClassDetailsByInviteCode(c.Request.Context(), inviteCode)
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.GlobalErrorResponse{
			Error:   "NOT_FOUND",
			Message: "Class not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Retrieved successfully",
		"data":    res,
	})
}
