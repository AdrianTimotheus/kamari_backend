package helper

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

var (
	errParamNotFound     = errors.New("failed to get value from the given param key")
	errInvalidReturnType = errors.New("failed to get value from the given param key")
)

func GetParam(c *gin.Context, key string, returnType string) (interface{}, error) {
	p := c.Param(key)
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
