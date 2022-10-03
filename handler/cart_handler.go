package handler

import (
	"github.com/gin-gonic/gin"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
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

func (h Handler) GetCartItems(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	response, err := h.cartService.GetCartItems(user)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h Handler) CreateCartItem(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	var req dto.AddCartItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperror.BadRequestError(ErrFailedBindingValidation.Error(), "INVALID_PAYLOAD"))
		return
	}

	err := h.cartService.CreateCartItem(user, req)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item added to cart",
	})
}

func (h Handler) UpdateCartItem(c *gin.Context) {
	var req dto.UpdateCartItemRequest

	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperror.BadRequestError(ErrFailedBindingValidation.Error(), "INVALID_PAYLOAD"))
		return
	}

	id, err := helper.GetParam(c, "id", "int")
	if err != nil {
		_ = c.Error(httperror.BadRequestError(helper.ErrInvalidReturnType.Error(), "BAD_PARAMS"))
		return
	}

	err = h.cartService.UpdateCartItem(id.(int), user.Id, req)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart updated",
	})
}

func (h Handler) DeleteCartItem(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	id, err := helper.GetParam(c, "id", "int")
	if err != nil {
		_ = c.Error(httperror.BadRequestError(helper.ErrInvalidReturnType.Error(), "BAD_PARAMS"))
		return
	}

	err = h.cartService.DeleteCartItem(id.(int), user.Id)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart item deleted successfully",
	})
}
