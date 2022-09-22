package handler

import (
	"github.com/gin-gonic/gin"
	"kamari.com/helper"
	"net/http"
)

func (h Handler) GetSellerDetail(c *gin.Context) {
	res, err := helper.GetParam(c, "username", "string")
	if err != nil {
		_ = c.Error(err)
		return
	}

	username := res.(string)
	merchant, err := h.merchantService.GetMerchantDetail(username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, merchant)
}
