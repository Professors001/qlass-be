package api

import (
	"qlass-be/adapters/api/rest"
)

type Handler struct {
	UserHandler *rest.UserHandler
}

func ProvideHandler(userHandler *rest.UserHandler) *Handler {
	return &Handler{
		UserHandler: userHandler,
	}
}
