package httperror

import (
	"net/http"
)

type AppError struct {
	StatusCode int    `json:"statusCode"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (err AppError) Error() string {
	return err.Message
}

func BadRequestError(message string, code string) AppError {
	if code == "" {
		code = "BAD_REQUEST"
	}
	return AppError{
		Code:       code,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func InternalServerError(message string) AppError {
	return AppError{
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

func UnauthorizedError() AppError {
	return AppError{
		Code:       "UNAUTHORIZED_ERROR",
		Message:    "Unauthorized error",
		StatusCode: http.StatusUnauthorized,
	}
}
