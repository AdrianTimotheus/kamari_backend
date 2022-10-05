package dto

type GetMerchantProductsResponse struct {
	Page      int       `json:"page"`
	TotalPage int       `json:"total_page"`
	Products  []Product `json:"products"`
}

type Product struct {
	ID           uint64  `json:"id"`
	Title        string  `json:"title"`
	Price        uint64  `json:"price"`
	ThumbnailURL string  `json:"thumbnail_url"`
	ProvinceName string  `json:"province"`
	CityName     string  `json:"city_name"`
	MerchantName string  `json:"merchant"`
	ItemSold     uint64  `json:"item_sold"`
	ReviewScore  float32 `json:"review_score"`
	ReviewCount  int     `json:"review_count"`
}
