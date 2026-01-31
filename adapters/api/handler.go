package api

import (
	"qlass-be/adapters/api/rest"
)

type Handler struct {
	UserHandler  *rest.UserHandler
	ClassHandler *rest.ClassHandler
}

func ProvideHandler(userHandler *rest.UserHandler, classHandler *rest.ClassHandler) *Handler {
	return &Handler{
		UserHandler:  userHandler,
		ClassHandler: classHandler,
	}
}
