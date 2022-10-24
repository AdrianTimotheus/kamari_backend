package dto

type CheckoutItemResponse struct {
	Merchant MerchantDetailResponse `json:"merchant"`
	Products []CheckoutProduct      `json:"products"`
	Couriers []SupportedCourier     `json:"couriers,omitempty"`
}

type CheckoutProduct struct {
	ProductId    uint64      `json:"product_id"`
	VariantID    uint64      `json:"variant_id"`
	ProductName  string      `json:"product_name"`
	RealPrice    uint64      `json:"real_price"`
	Price        uint64      `json:"price"`
	TotalPrice   uint64      `json:"total-price"`
	ThumbnailURL string      `json:"thumbnail_url"`
	MainVariant  MainVariant `json:"main_variant,omitempty"`
	SubVariant   SubVariant  `json:"sub_variant,omitempty"`
	Quantity     int         `json:"quantity"`
	ItemWeight   uint64      `json:"item_weight"`
	TotalWeight  uint64      `json:"total_weight"`
}

type MainVariant struct {
	Name string `json:"name,omitempty"`
	Item string `json:"item,omitempty"`
}

type SubVariant struct {
	Name string `json:"name,omitempty"`
	Item string `json:"item,omitempty"`
}

type SupportedCourier struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}
