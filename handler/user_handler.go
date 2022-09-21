package handler

import (
	"github.com/gin-gonic/gin"
	"kamari.com/dto"
	"net/http"
)

func (h Handler) Register(c *gin.Context) {
	var (
		req *dto.RegisterRequest
		err error
	)

	if err = c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(ErrFailedBindingValidation)
		return
	}

	err = h.userService.CreateNewUser(req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
	})
}
