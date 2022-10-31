package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
	"kamari.com/model"
)

func (h Handler) ProceedOrder(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	var payload dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(httperror.BadRequestError("Missing Required Field", "MISSING_REQUIRED_FIELD"))
		return
	}
	fmt.Printf("%+v\n", payload)
	fmt.Printf("%+v\n", user)
	resp, err := h.orderService.CreateOrder(user, payload)
	if err != nil {
		_ = c.Error(err)
	}
	c.JSON(http.StatusCreated, resp)
}

func (h Handler) ProceedOrderSLP(c *gin.Context) {
	val := c.Query("s")

	var req *dto.SealabsPayRes

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperror.BadRequestError(ErrFailedBindingValidation.Error(), "INVALID_PAYLOAD"))
		return
	}

	if req.Status == "TXN_FAILED" {
		_, err := h.transactionService.CancelOrder(val)
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": req.Status,
		})
		return
	}

	resp, err := h.orderService.CreateOrderSLP(val)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h Handler) GetMerchantOrders(c *gin.Context) {
	var query struct {
		Keyword string `form:"keyword"`
		helper.Filters
	}

	err := c.BindQuery(&query)
	if err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_QUERY_PARAMS"))
		return
	}

	user, err := helper.GetUserContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httperror.UnauthorizedError())
		return
	}

	merchant, err := h.merchantService.GetMerchantDetail(user.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	orders, metadata, err := h.orderService.GetMerchantOrder(merchant.Id, query.Keyword, &query.Filters)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"__metadata": metadata,
		"orders":     orders,
	})
}

func (h Handler) GetMerchantOrderDetail(c *gin.Context) {
	i, err := helper.GetParam(c, "id", "uint64")
	if err != nil {
		_ = c.Error(httperror.BadRequestError(helper.ErrInvalidReturnType.Error(), "BAD_PARAMS"))
		return
	}
	orderId := i.(uint64)

	user, err := helper.GetUserContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httperror.UnauthorizedError())
		return
	}

	merchant, err := h.merchantService.GetMerchantDetail(user.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	order, err := h.orderService.GetOrderDetails(orderId, merchant.Id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order": order,
	})

}
func (h Handler) GetMerchantOrderSummary(c *gin.Context) {
	user, err := helper.GetUserContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httperror.UnauthorizedError())
		return
	}

	merchant, err := h.merchantService.GetMerchantDetail(user.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	orders, err := h.orderService.GetMerchantOrderSummary(merchant.Id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}

func (h Handler) UpdateOrderStatus(c *gin.Context) {
	id, err := helper.GetParam(c, "id", "uint64")
	if err != nil {
		_ = c.Error(httperror.BadRequestError(helper.ErrInvalidReturnType.Error(), "BAD_PARAMS"))
		return
	}

	userVal, _ := c.Get("user")
	user := userVal.(model.User)
	dtoResp, err := h.orderService.UpwardOrderStatus(id.(uint64), user)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dtoResp)

}
