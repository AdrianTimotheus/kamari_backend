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
	ErrParamNotFound     = errors.New("failed to get value from the given param key")
	ErrInvalidReturnType = errors.New("failed to get value from the given param key")
)

func GetParam(c *gin.Context, key string, returnType string) (interface{}, error) {
	p := c.Param(key)
	p = strings.TrimSpace(p)
	if p == "" {
		return nil, ErrParamNotFound
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
		return nil, ErrInvalidReturnType
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
		if int(product.ReviewScore) >= rating {
			res = append(res, product)
		}
	}

	return res
}

func FilterReviewsByMedia(reviews []model.ProductReview, metadata *Metadata) ([]model.ProductReview, *Metadata) {
	var res []model.ProductReview

	for _, review := range reviews {
		if len(review.ReviewPhotos) > 0 {
			res = append(res, review)
		}
	}

	*metadata = CalculateMetaData(len(res), metadata.CurrentPage, metadata.Limit)
	return res, metadata
}

func CountPriceRange(variants []model.ProductVariant) []uint64 {
	var (
		res = make([]uint64, 2)
		min uint64
		max uint64
	)

	min, max = variants[0].Price, variants[0].Price
	for _, variant := range variants {
		if variant.Price < min {
			min = variant.Price
		}
		if variant.Price > max {
			max = variant.Price
		}
	}

	res[0], res[1] = min, max
	return res
}
