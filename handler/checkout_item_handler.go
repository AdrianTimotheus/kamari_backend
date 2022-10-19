package handler

import (
	"github.com/gin-gonic/gin"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
	"net/http"
)

func (h Handler) AddCheckOutItems(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)
	var payloads []dto.PostCheckoutRequest
	if err := c.ShouldBindJSON(&payloads); err != nil {
		_ = c.Error(httperror.BadRequestError("missing required field", "MISSING_REQUIRED_FIELD"))
		return
	}
	resp, err := h.checkoutItemService.CreateCheckoutItem(user, payloads)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h Handler) GetCheckOutItems(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)
	resp, err := h.checkoutItemService.GetCheckoutItemByUser(user)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
