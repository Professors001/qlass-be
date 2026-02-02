package api

import (
	"qlass-be/adapters/api/rest"
)

type Handler struct {
	UserHandler       *rest.UserHandler
	ClassHandler      *rest.ClassHandler
	AttachmentHandler *rest.AttachmentHandler
}

func ProvideHandler(userHandler *rest.UserHandler, classHandler *rest.ClassHandler, attachmentHandler *rest.AttachmentHandler) *Handler {
	return &Handler{
		UserHandler:       userHandler,
		ClassHandler:      classHandler,
		AttachmentHandler: attachmentHandler,
	}
}
