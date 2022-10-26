package dto

type CreateOrderRequest struct {
	//PaymentMethod string          `json:"payment_method" binding:"required"`
	//CardNumber    string          `json:"card_number"`
	TransactionId uint64          `json:"transaction_id" binding:"required"`
	Orders        []Order         `json:"orders" binding:"required"`
	Products      []ItemToProceed `json:"products" binding:"required"`
}

type ItemToProceed struct {
	ProductId  uint64 `json:"product_id" binding:"required"`
	VariantID  uint64 `json:"variant_id" binding:"required"`
	RealPrice  uint64 `json:"real_price" binding:"required"`
	Price      uint64 `json:"price" binding:"required"`
	TotalPrice uint64 `json:"total_price" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required"`
	MerchantId uint64 `json:"merchant_id" binding:"required"`
}

type Order struct {
	UserTotalPayment         uint64  `json:"user_total_payment" binding:"required"`
	CourierId                uint64  `json:"courier_id" binding:"required"`
	ShippingFee              uint64  `json:"shipping_fee" binding:"required"`
	MerchantVoucherId        *uint64 `json:"merchant_voucher_id,omitempty"`
	MarketplaceVoucherId     *uint64 `json:"marketplace_voucher_id,omitempty"`
	DstAddressId             uint64  `json:"dst_address_id" binding:"required"`
	MerchantId               uint64  `json:"merchant_id" binding:"required"`
	MerchantIncome           uint64  `json:"merchant_income" binding:"required"`
	TotalProductsPrice       uint64  `json:"total_products_price" binding:"required"`
	MerchantTotalDiscount    uint64  `json:"merchant_total_discount"`
	MarketPlaceTotalDiscount uint64  `json:"market_place_total_discount"`
}
