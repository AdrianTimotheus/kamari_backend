package dto

import "kamari.com/model"

type GetSellerAdminProductsResponse struct {
	Page      int                  `json:"page"`
	TotalPage int                  `json:"total_page"`
	Products  []SellerAdminProduct `json:"products"`
}

type SellerAdminProduct struct {
	ID              uint64                 `json:"id"`
	Title           string                 `json:"title"`
	ThumbnailURL    string                 `json:"thumbnail_url"`
	ItemSold        uint64                 `json:"item_sold"`
	ReviewScore     float32                `json:"review_score"`
	ReviewCount     int                    `json:"review_count"`
	ProductVariants []model.ProductVariant `json:"product_variants"`
}
