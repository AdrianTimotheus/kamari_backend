package handler

import (
	"github.com/gin-gonic/gin"
	"kamari.com/model"
	"net/http"
)

func (h Handler) GetWalletDetail(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	response, err := h.walletService.GetUserWalletDetail(user)
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, response)
}
