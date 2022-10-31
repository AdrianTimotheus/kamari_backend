package handler

import (
	"fmt"
	"kamari.com/dto"
	"kamari.com/model"
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

func (h Handler) AddNewProduct(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)
	var payload dto.AddProductReq
	fmt.Println(payload)

	if err := c.ShouldBind(&payload); err != nil {
		//_ = c.Error(httperror.BadRequestError("Invalid request body", "INVALID_REQ_BODY"))
		c.Error(err)
		return
	}
	fmt.Println(payload)

	err := h.productService.AddNewProduct(payload, user)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product successfully created",
	})
}

func (h Handler) EditProduct(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)
	var payload dto.EditProductReq

	if err := c.ShouldBind(&payload); err != nil {
		_ = c.Error(httperror.BadRequestError("Invalid query", "INVALID_QUERY"))
		return
	}

	id, err := helper.GetParam(c, "id", "uint64")
	if err != nil {
		_ = c.Error(httperror.BadRequestError(helper.ErrInvalidReturnType.Error(), "BAD_PARAMS"))
		return
	}

	err = h.productService.EditProduct(payload, user, id.(uint64))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product successfully updated",
	})
}

func (h Handler) DeleteProduct(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)
	id, err := helper.GetParam(c, "id", "uint64")
	if err != nil {
		_ = c.Error(httperror.BadRequestError(helper.ErrInvalidReturnType.Error(), "BAD_PARAMS"))
		return
	}

	err = h.productService.DeleteProduct(id.(uint64), user)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product successfully deleted",
	})
}
