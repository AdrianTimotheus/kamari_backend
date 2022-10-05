package handler

import (
	"kamari.com/dto"
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

func (h Handler) GetProductDetail(c *gin.Context) {
	productId, err := helper.GetParam(c, "id", "int")
	if err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_PARAM"))
		return
	}

	product, err := h.productService.GetProductDetail(uint64(productId.(int)))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h Handler) GetSimilarProducts(c *gin.Context) {
	var payload dto.GetSimilarProductRequest

	if err := c.ShouldBind(&payload); err != nil {
		_ = c.Error(httperror.BadRequestError("Invalid query", "INVALID_QUERY"))
		return
	}
	resp, err := h.productService.GetSimilarProduct(payload)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
