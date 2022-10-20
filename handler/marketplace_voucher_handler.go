package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h Handler) GetAllMarketplaceActiveVoucher(c *gin.Context) {
	resp, err := h.marketplaceVoucherService.GetAllActiveMarketplaceVoucher()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
