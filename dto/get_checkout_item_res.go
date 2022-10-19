package dto

type CheckoutItemResponse struct {
	Merchant MerchantDetailResponse
	Products []CheckoutProduct
}

type CheckoutProduct struct {
	ProductId    uint64      `json:"product_id"`
	VariantID    uint64      `json:"variant_id"`
	ProductName  string      `json:"product_name"`
	Price        uint64      `json:"price"`
	ThumbnailURL string      `json:"thumbnail_url"`
	MainVariant  MainVariant `json:"main_variant,omitempty"`
	SubVariant   SubVariant  `json:"sub_variant,omitempty"`
	Quantity     int         `json:"quantity"`
}

type MainVariant struct {
	Name string `json:"name,omitempty"`
	Item string `json:"item,omitempty"`
}

type SubVariant struct {
	Name string `json:"name,omitempty"`
	Item string `json:"item,omitempty"`
}
