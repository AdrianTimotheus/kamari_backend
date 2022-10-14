package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/api/oauth2/v2"
	"kamari.com/config"
	"kamari.com/httperror"
	"kamari.com/model"
	"net/http"
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

func VerifyGoogleToken(idToken string) (*oauth2.Tokeninfo, error) {
	var httpClient = &http.Client{}
	oauth2Service, err := oauth2.New(httpClient)
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return nil, err
	}
	return tokenInfo, nil
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

func AuthorizeStepUpJWT(c *gin.Context) {
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
	scopeJson, err := json.Marshal(claims["Scope"])
	scope := string(scopeJson[:])
	c.Set("user", user)
	c.Set("Scope", scope)
}

func GetUserFromToken(encodedToken string) (*model.User, error) {
	unauthorizedErr := httperror.UnauthorizedError()
	token, err := validateToken(encodedToken)
	if err != nil || !token.Valid {
		return nil, unauthorizedErr
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	fmt.Println(claims)
	if !ok {
		return nil, unauthorizedErr
	}
	userJson, err := json.Marshal(claims["user"])
	var user model.User
	err = json.Unmarshal(userJson, &user)
	if err != nil {
		return nil, unauthorizedErr
	}
	return &user, nil
}
