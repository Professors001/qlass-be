package rest

import (
	"net/http"
	"qlass-be/dtos"
	"qlass-be/infrastructure/middleware"
	"qlass-be/usecases"

	"github.com/gin-gonic/gin"
)

type AttachmentHandler struct {
	UseCase usecases.AttachmentUseCase
}

func NewAttachmentHandler(attachmentUseCase usecases.AttachmentUseCase) *AttachmentHandler {
	return &AttachmentHandler{
		UseCase: attachmentUseCase,
	}
}

func (h *AttachmentHandler) UploadAttachment(ctx *gin.Context) {
	val, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dtos.GlobalErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "User context not found",
		})
		return
	}

	claims, ok := val.(*middleware.JWTCustomClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: "Failed to cast user context",
		})
		return
	}

	fileHeader, err := ctx.FormFile("attachment")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.GlobalErrorResponse{
			Error:   "BAD_REQUEST",
			Message: "No attachment uploaded",
		})
		return
	}

	attachment, err := h.UseCase.UploadAttachment(claims.UserId, fileHeader)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.GlobalErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: "Failed to upload attachment: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.UploadAttachmentResponseDto{
		AttachmentID: attachment.ID,
		FileURL:      attachment.FileURL,
	})

}

func (h *AttachmentHandler) GetAttachment() {

}

func (h *AttachmentHandler) GetAttachmentsByCourseMaterialID() {

}

func (h *AttachmentHandler) GetAttachmentsBySubmissionID() {

}

func (h *AttachmentHandler) UpDateAttachment() {

}

func (h *AttachmentHandler) DeleteAttachment() {

}
