package middleware

import (
	"github.com/gin-gonic/gin"
	"kamari.com/httperror"
)

func ErrorHandler(c *gin.Context) {
	if len(c.Errors) == 0 {
		return
	}

	firstError := c.Errors[0].Err
	appError, ok := firstError.(httperror.AppError)
	if ok {
		c.JSON(appError.StatusCode, appError)
		return
	}

	serverErr := httperror.InternalServerError(firstError.Error())
	c.JSON(serverErr.StatusCode, serverErr)
}
