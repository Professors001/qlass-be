package dtos

type GlobalErrorResponse struct {
	Error string `json:"error"`
	Message string `json:"message"`
}