package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
	"net/http"
)

func (h Handler) ProceedOrder(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	var payload dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(httperror.BadRequestError("Missing Required Field", "MISSING_REQUIRED_FIELD"))
		return
	}
	fmt.Printf("%+v\n", payload)
	fmt.Printf("%+v\n", user)
	resp, err := h.orderService.CreateOrder(user, payload)
	if err != nil {
		_ = c.Error(err)
	}
	c.JSON(http.StatusCreated, resp)
}
