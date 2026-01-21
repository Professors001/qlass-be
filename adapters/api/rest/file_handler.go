package rest

import (
	"net/http"
	"qlass-be/config"
	"qlass-be/usecases"

	"github.com/gin-gonic/gin"
)

type FileController struct {
	fileUseCase usecases.FileUseCase
	cfg         *config.Config
}

func NewFileController(fileUseCase usecases.FileUseCase, cfg *config.Config) *FileController {
	return &FileController{
		fileUseCase: fileUseCase,
		cfg:         cfg,
	}
}

func (c *FileController) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	objectName, err := c.fileUseCase.UploadFile(ctx.Request.Context(), file, c.cfg.MinioBucketName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":     "File uploaded successfully",
		"object_name": objectName,
	})
}

func (c *FileController) GetUrl(ctx *gin.Context) {
	objectName := ctx.Query("object_name")
	if objectName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "object_name query parameter is required"})
		return
	}

	url, err := c.fileUseCase.GetFileUrl(ctx.Request.Context(), objectName, c.cfg.MinioBucketName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate URL", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"object_name": objectName,
		"url":         url,
	})
}
