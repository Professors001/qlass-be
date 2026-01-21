package api

import (
	"qlass-be/adapters/api/rest"
)

type Handler struct {
	UserHandler  *rest.UserHandler
	ClassHandler *rest.ClassHandler
	FileHandler  *rest.FileController
}

func ProvideHandler(userHandler *rest.UserHandler, classHandler *rest.ClassHandler, fileHandler *rest.FileController) *Handler {
	return &Handler{
		UserHandler:  userHandler,
		ClassHandler: classHandler,
		FileHandler:  fileHandler,
	}
}
