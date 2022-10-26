package handler

import (
	"github.com/gin-gonic/gin"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
	"net/http"
)

func (h Handler) CreateProductPromotion(c *gin.Context) {
	var payload dto.ProductPromoRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_PAYLOAD"))
		return
	}

	user, err := helper.GetUserContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httperror.UnauthorizedError())
		return
	}

	err = h.productPromoService.CreateProductPromo(&payload, user.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully created voucher",
	})
}

func (h Handler) GetProductPromos(c *gin.Context) {
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

	u, err := h.merchantService.GetMerchantDetail(user.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	promos, metadata, err := h.productPromoService.GetProductPromos(u.Id, query.Keyword, &query.Filters)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"__metadata":    metadata,
		"product_promo": promos,
	})
}

func (h Handler) GetProductPromoDetail(c *gin.Context) {
	p, err := helper.GetParam(c, "id", "uint64")
	if err != nil {
		_ = httperror.BadRequestError(err.Error(), "BAD_PARAMS")
		return
	}

	promoId := p.(uint64)
	voucher, err := h.productPromoService.GetProductPromoDetail(promoId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, voucher)
}

func (h Handler) DeletePromoItems(c *gin.Context) {
	p, err := helper.GetParam(c, "id", "uint64")
	if err != nil {
		_ = httperror.BadRequestError(ErrFailedBindingValidation.Error(), "BAD_PARAMS")
		return
	}

	promoId := p.(uint64)

	promoItems := struct {
		PromoItemsId []uint64 `json:"promo_items_id"`
	}{}
	if err = c.ShouldBindJSON(&promoItems); err != nil {
		_ = c.Error(httperror.BadRequestError(ErrFailedBindingValidation.Error(), "INVALID_PAYLOAD"))
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

	if err = h.productPromoService.DeletePromoItem(merchant.Id, promoId, promoItems.PromoItemsId); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "promo items successfully deleted",
	})
}

func (h Handler) EditProductPromo(c *gin.Context) {
	p, err := helper.GetParam(c, "id", "uint64")
	if err != nil {
		_ = httperror.BadRequestError(err.Error(), "BAD_PARAMS")
		return
	}

	promoId := p.(uint64)

	var payload dto.ProductPromoUpdate
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_PAYLOAD"))
		return
	}

	user, err := helper.GetUserContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httperror.UnauthorizedError())
		return
	}

	err = h.productPromoService.UpdateProductPromo(&payload, user.Username, promoId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully updated voucher",
	})
}

func (h Handler) InsertPromoItems(c *gin.Context) {
	p, err := helper.GetParam(c, "id", "uint64")

	if err != nil {
		_ = httperror.BadRequestError(err.Error(), "BAD_PARAMS")
		return
	}

	promoId := p.(uint64)

	type promoItems struct {
		PromotionItems []dto.PromoItems `json:"promotionItems" binding:"required"`
	}

	var payload promoItems

	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_PAYLOAD"))
		return

	}

	err = h.productPromoService.InsertPromoItems(payload.PromotionItems, promoId)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully insert new product promo",
	})

}

func (h Handler) GetAvailableMerchantProductForPromo(c *gin.Context) {
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

	u, err := h.merchantService.GetMerchantDetail(user.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	availProducts, metadata, err := h.productPromoService.GetAvailableProductForPromo(u.Id, query.Keyword, &query.Filters)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"__metadata": metadata,
		"products":   availProducts,
	})
}

func (h Handler) EndMerchantPromo(c *gin.Context) {
	p, err := helper.GetParam(c, "id", "uint64")
	if err != nil {
		_ = httperror.BadRequestError(err.Error(), "BAD_PARAMS")
		return
	}

	promoId := p.(uint64)
	err = h.productPromoService.EndProductPromo(promoId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully ended a promo",
	})
}
