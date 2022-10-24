package middleware

import (
	"github.com/gin-gonic/gin"
	"kamari.com/helper"
	"kamari.com/httperror"
	"kamari.com/model"
	"net/http"
)

func AuthorizeAdmin(c *gin.Context) {
	user, err := helper.GetUserContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httperror.UnauthorizedError())
		return
	}

	if err != nil || user.RoleId != model.RoleAdmin {
		c.AbortWithStatusJSON(http.StatusForbidden, httperror.ForbiddenError())
		return
	}
}
