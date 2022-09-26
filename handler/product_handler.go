package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h Handler) GetRecommendedProduct(c *gin.Context) {
	resp, err := h.productService.GetTopSellerProduct()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
