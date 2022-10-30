package dto

type OrderHighlightResponse struct {
	Id               uint64             `json:"id"`
	User             UserDetail         `json:"user" gorm:"-"`
	Invoice          string             `json:"invoice"`
	OrderStatus      string             `json:"order_status"`
	ShippingOption   ShippingOption     `json:"shipping_option" gorm:"-"`
	OrderItems       []HighlightProduct `json:"order_items" gorm:"-"`
	UserTotalPayment uint64             `json:"user_total_payment"`
	MerchantIncome   uint64             `json:"merchant_income"`
	ShippingFee      uint64             `json:"shipping_fee"`
	UserId           uint64             `json:"-"`
	DstAddressId     uint64             `json:"-"`
	CourierId        uint64             `json:"-"`
	OrderStatusId    uint64             `json:"-"`
}

type HighlightProduct struct {
	Id               uint64      `json:"id"`
	ProductId        uint64      `json:"product_id"`
	ProductName      string      `json:"product_name,omitempty"`
	ProductVariantId uint64      `json:"-"`
	RealPrice        uint64      `json:"real_price"`
	Price            uint64      `json:"price"`
	Quantity         int         `json:"quantity"`
	TotalPrice       uint64      `json:"total_price"`
	ThumbnailURL     string      `json:"thumbnail_url,omitempty"`
	MainVariant      MainVariant `json:"main_variant,omitempty" gorm:"embedded"`
	SubVariant       SubVariant  `json:"sub_variant,omitempty" gorm:"embedded"`
}

type ShippingOption struct {
	Id      uint64 `json:"id"`
	Courier string `json:"courier"`
}

type UserDetail struct {
	Name    string  `json:"name"`
	Phone   string  `json:"phone"`
	PicUrl  string  `json:"pic_url"`
	Address Address `json:"address,omitempty" gorm:"-"`
}

type Address struct {
	RecipientName string       `json:"recipient_name"`
	Phone         string       `json:"phone"`
	City          CityResponse `json:"city" gorm:"-"`
	AddressDetail string       `json:"address_detail" gorm:"column:detail_address"`
	CityId        uint64       `json:"-"`
}

type CityResponse struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}
