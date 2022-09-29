package handler

import (
	"github.com/gin-gonic/gin"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
	"net/http"
)

func (h Handler) Favorite(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	var payload dto.FavoriteRequest
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		_ = c.Error(httperror.BadRequestError("Missing Required Query Param", "MISSING_REQUIRED_QUERY_PARAM"))
		return
	}

	response, err := h.wishlistService.Favorite(user, payload.ProductId)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}
