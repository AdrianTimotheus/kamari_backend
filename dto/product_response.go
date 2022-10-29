package dto

import "kamari.com/model"

type ProductResponse struct {
	Id               uint64               `json:"id"`
	Title            string               `json:"title"`
	Price            uint64               `json:"price"`
	ItemSold         uint64               `json:"item_sold"`
	City             string               `json:"city"`
	ThumbnailUrl     string               `json:"thumbnail_url"`
	FavouriteCount   uint64               `json:"favourite_count"`
	Availability     bool                 `json:"availability"`
	MerchantId       uint64               `json:"merchant_id"`
	MerchantUsername string               `json:"merchant_username"`
	ReviewScore      float32              `json:"review_score"`
	ReviewCount      int                  `json:"review_count"`
	CategoryId       uint64               `json:"category_id"`
	Category         model.Category       `json:"category"`
	SellerCategory   model.SellerCategory `json:"seller_category"`
}
