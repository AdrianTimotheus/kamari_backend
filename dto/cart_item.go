package dto

type CartItem struct {
	Id              uint64 `json:"id"`
	Name            string `json:"name"`
	MainVariant     string `json:"main_variant"`
	MainVariantItem string `json:"main_variant_item"`
	SubVariant      string `json:"sub_variant"`
	SubVariantItem  string `json:"sub_variant_item"`
	PicUrl          string `json:"pic_url"`
	Price           uint64 `json:"price"`
	Quantity        uint64 `json:"quantity"`
}

type CartItemRaw struct {
	Id              uint64 `json:"id"`
	Name            string `json:"name"`
	MainVariant     string `json:"main_variant"`
	MainVariantItem string `json:"main_variant_item"`
	SubVariant      string `json:"sub_variant"`
	SubVariantItem  string `json:"sub_variant_item"`
	PicUrl          string `json:"pic_url"`
	Price           uint64 `json:"price" column:"price"`
	Quantity        uint64 `json:"quantity"`
	MerchantId      uint64 `json:"merchant_id" column:"merchant_id"`
	MerchantName    string `json:"merchant_name" column:"merchant_name"`
	MerchantPicUrl  string `json:"merchant_pic_url" column:"pic_url"`
}
