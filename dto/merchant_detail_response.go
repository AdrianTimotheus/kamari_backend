package dto

import "time"

type MerchantDetailResponse struct {
	Id                uint64    `json:"id"`
	Name              string    `json:"name"`
	Username          string    `json:"username"`
	Description       string    `json:"description"`
	ProfilePicture    string    `json:"profile_picture"`
	FollowerCount     uint64    `json:"follower_count"`
	ProductCount      uint64    `json:"product_count"`
	ShopLocation      string    `json:"shop_location"`
	MerchantCreatedAt time.Time `json:"merchant_created_at"`
}
