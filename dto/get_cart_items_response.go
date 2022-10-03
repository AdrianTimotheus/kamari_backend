package dto

type CartItemsResponse struct {
	Merchant CartItemMerchant `json:"merchant"`
	Items    []CartItem       `json:"items"`
}
