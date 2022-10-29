package model

type SellerOrderDetail struct {
	Id              uint64 `json:"id"`
	ProductName     string `json:"product_name"`
	ThumbnailUrl    string `json:"thumbnail_url"`
	ProductId       uint64 `json:"product_id"`
	MainVariant     string `json:"main_variant"`
	MainVariantItem string `json:"main_variant_item"`
	SubVariant      string `json:"sub_variant"`
	SubVariantItem  string `json:"sub_variant_item"`
	Quantity        uint64 `json:"quantity"`
}
