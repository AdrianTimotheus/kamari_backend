package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"kamari.com/config"
	"kamari.com/httperror"
	"kamari.com/model"
	"strings"
)

func validateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, httperror.UnauthorizedError()
		}
		return config.Config.JWTSecret, nil
	})
}

func AuthorizeJWT(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	fmt.Println("auth header", authHeader)

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
	c.Set("user", user)
}
