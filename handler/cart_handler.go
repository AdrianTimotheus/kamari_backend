package handler

import (
	"github.com/gin-gonic/gin"
	"kamari.com/model"
	"net/http"
)

func (h Handler) GetCartHoverItems(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	response, err := h.cartService.GetCartHoverItems(user)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}
