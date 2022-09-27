package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"kamari.com/helper"
	"kamari.com/httperror"
)

func (h Handler) GetProducts(c *gin.Context) {
	var query struct {
		Keyword string `form:"keyword"`
		helper.Filters
	}

	err := c.BindQuery(&query)
	if err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_QUERY_PARAMS"))
		return
	}

	products, metadata, err := h.productService.GetProducts(query.Keyword, &query.Filters)
	c.JSON(http.StatusOK, gin.H{
		"__metadata": metadata,
		"products":   products,
	})
}

func (h Handler) GetRecommendedProduct(c *gin.Context) {
	resp, err := h.productService.GetTopSellerProduct()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
