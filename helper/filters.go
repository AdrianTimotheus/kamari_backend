package helper

import (
	"math"
)

type Filters struct {
	Page          int     `form:"page"`
	Limit         int     `form:"limit"`
	SortBy        string  `form:"sortBy"`
	Sort          string  `form:"sort"`
	SortDirection string  `form:"sortDirection"`
	Location      string  `form:"location"`
	FilterBy      string  `form:"filterBy"`
	MaxPrice      *uint64 `form:"maxPrice"`
	MinPrice      *uint64 `form:"minPrice"`
	Rating        *int    `form:"rating"`
	CategoryId    *uint64 `form:"categoryId"`
	SubCategoryId *uint64 `form:"subCategoryId"`
	VoucherStatus *string `form:"voucherStatus"`
	PromoStatus   *string `form:"promoStatus"`
	OrderStatus   *string `form:"orderStatus"`
}

func (f Filters) Offset() int {
	return (f.Page - 1) * f.Limit
}

type categoryMetadata struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

func NewCategoryMetadata(id uint64, name string) categoryMetadata {
	return categoryMetadata{
		Id:   id,
		Name: name,
	}
}

type Metadata struct {
	CurrentPage  int                `json:"currentPage,omitempty"`
	Limit        int                `json:"limit,omitempty"`
	FirstPage    int                `json:"firstPage,omitempty"`
	LastPage     int                `json:"lastPage,omitempty"`
	TotalRecords int                `json:"totalRecords,omitempty"`
	Cities       []string           `json:"cities,omitempty"`
	Categories   []categoryMetadata `json:"categories,omitempty"`
}

func CalculateMetaData(totalRecords, page, limit int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		Limit:        limit,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(limit))),
		TotalRecords: totalRecords,
	}
}

func ReadString(val, defaultValue string) string {
	if val == "" {
		return defaultValue
	}

	return val
}

func ReadInt(val, defaultValue int) int {
	if val == 0 {
		return defaultValue
	}

	return val
}
