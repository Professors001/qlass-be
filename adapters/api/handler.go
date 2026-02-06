package api

import (
	"qlass-be/adapters/api/rest"
)

type Handler struct {
	UserHandler          *rest.UserHandler
	ClassHandler         *rest.ClassHandler
	AttachmentHandler    *rest.AttachmentHandler
	ClassMaterialHandler *rest.MaterialHandler
	SubmissionHandler    *rest.SubmissionHandler
}

func ProvideHandler(
	userHandler *rest.UserHandler,
	classHandler *rest.ClassHandler,
	attachmentHandler *rest.AttachmentHandler,
	classMaterialHandler *rest.MaterialHandler,
	submissionHandler *rest.SubmissionHandler,
) *Handler {
	return &Handler{
		UserHandler:          userHandler,
		ClassHandler:         classHandler,
		AttachmentHandler:    attachmentHandler,
		ClassMaterialHandler: classMaterialHandler,
		SubmissionHandler:    submissionHandler,
	}
}
