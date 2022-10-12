package handler

import (
	"github.com/gin-gonic/gin"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
	"kamari.com/model"
	"net/http"
)

func (h Handler) GetAddress(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	response, err := h.addressService.GetAddress(user.Id)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h Handler) EditAddress(c *gin.Context) {
	var req dto.EditAddressReq

	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperror.BadRequestError(ErrFailedBindingValidation.Error(), "INVALID_PAYLOAD"))
		return
	}

	response, err := h.addressService.EditAddress(user.Id, req)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h Handler) AddNewAddress(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	var req dto.AddNewAddressRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperror.BadRequestError(ErrFailedBindingValidation.Error(), "INVALID_PAYLOAD"))
		return
	}

	response, err := h.addressService.AddNewAddress(user.Id, req)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h Handler) DeleteAddress(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	id, err := helper.GetParam(c, "id", "int")
	if err != nil {
		_ = c.Error(httperror.BadRequestError(helper.ErrInvalidReturnType.Error(), "BAD_PARAMS"))
		return
	}

	err = h.addressService.DeleteAddress(user.Id, id.(int))

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "address successfully deleted",
	})
}

func (h Handler) GetDefaultAddress(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	response, err := h.addressService.GetDefaultAddress(user.Id)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h Handler) SetDefaultAddress(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	id, err := helper.GetParam(c, "id", "int")
	if err != nil {
		_ = c.Error(httperror.BadRequestError(helper.ErrInvalidReturnType.Error(), "BAD_PARAMS"))
		return
	}

	response, err := h.addressService.SetDefaultAddress(user.Id, id.(int))

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}
