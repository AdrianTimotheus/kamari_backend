package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h Handler) GetMarketplacePromoBanner(c *gin.Context) {
	response, err := h.marketplaceService.GetPromoBanner()

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}
