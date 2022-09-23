package helper

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	errParamNotFound     = errors.New("failed to get value from the given param key")
	errInvalidReturnType = errors.New("failed to get value from the given param key")
)

func GetParam(c *gin.Context, key string, returnType string) (interface{}, error) {
	p := c.Param(key)
	p = strings.TrimSpace(p)
	if p == "" {
		return nil, errParamNotFound
	}

	switch returnType {
	case "string":
		return p, nil
	case "int":
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		return n, nil
	default:
		return nil, errInvalidReturnType
	}
}
