package helper

import (
	"errors"
	"strconv"
	"strings"

	"kamari.com/dto"
	"kamari.com/model"

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

func GetMinProductVarPrice(s []model.ProductVariant) uint64 {
	if len(s) < 1 {
		return 0
	}

	min := s[0].Price
	for i := 1; i < len(s); i++ {
		if s[i].Price < min {
			min = s[i].Price
		}
	}

	return min
}

func CalcProductRating(pv []model.OrderDetail) float32 {
	if len(pv) < 1 {
		return 0
	}

	var sum float32
	for _, o := range pv {
		sum += o.ProductReview.Score
	}

	return sum / float32(len(pv))
}

func readCSV(s string) []string {
	parts := strings.Split(s, ",")
	return parts
}

func inSlice(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}

	return false
}

func FilterBySellerLocation(products []dto.ProductResponse, location string) []dto.ProductResponse {
	var (
		loc = readCSV(location)
		res []dto.ProductResponse
	)

	for _, product := range products {
		if inSlice(loc, product.City) {
			res = append(res, product)
		}
	}

	return res
}

func FilterProductPrice(products []dto.ProductResponse, priceType string, price uint64) []dto.ProductResponse {
	var res []dto.ProductResponse

	for _, product := range products {
		switch priceType {
		case "minPrice":
			if product.Price >= price {
				res = append(res, product)
			}
		case "maxPrice":
			if product.Price <= price {
				res = append(res, product)
			}
		}
	}
	return res
}

func FilterProductByRating(products []dto.ProductResponse, rating int) []dto.ProductResponse {
	var res []dto.ProductResponse

	for _, product := range products {
		if product.ReviewScore >= rating {
			res = append(res, product)
		}
	}

	return res
}
