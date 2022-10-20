package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"kamari.com/helper"
	"kamari.com/httperror"
	"kamari.com/model"
	"strings"
)

func AuthorizeStepUpJWT(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	s := strings.Split(authHeader, "Bearer ")
	unauthorizedErr := httperror.UnauthorizedError()
	if len(s) < 2 {
		c.AbortWithStatusJSON(unauthorizedErr.StatusCode, unauthorizedErr)
		return
	}
	encodedToken := s[1]
	token, err := validateToken(encodedToken)
	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(unauthorizedErr.StatusCode, unauthorizedErr)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	fmt.Println(claims)
	if !ok {
		c.AbortWithStatusJSON(unauthorizedErr.StatusCode, unauthorizedErr)
		return
	}

	userJson, err := json.Marshal(claims["user"])
	var user model.User
	err = json.Unmarshal(userJson, &user)
	if err != nil {
		c.AbortWithStatusJSON(unauthorizedErr.StatusCode, unauthorizedErr)
		return
	}
	scopeJson, err := json.Marshal(claims["Scope"])
	scope := string(scopeJson[:])

	redis := helper.NewRedisClient(1)
	stepUpToken, err := redis.GetData("stepUp-" + user.Email)
	if err != nil {
		c.AbortWithStatusJSON(unauthorizedErr.StatusCode, unauthorizedErr)
		return
	}
	if stepUpToken != encodedToken {
		c.AbortWithStatusJSON(unauthorizedErr.StatusCode, unauthorizedErr)
		return
	}
	c.Set("user", user)
	c.Set("Scope", scope)
}
