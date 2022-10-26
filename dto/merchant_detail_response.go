package dto

import "time"

type MerchantDetailResponse struct {
	Id                uint64    `json:"id"`
	Name              string    `json:"name"`
	Username          string    `json:"username,omitempty"`
	Description       string    `json:"description,omitempty"`
	ProfilePicture    string    `json:"profile_picture,omitempty"`
	FollowerCount     uint64    `json:"follower_count,omitempty"`
	ProductCount      uint64    `json:"product_count,omitempty"`
	ShopLocation      string    `json:"shop_location,omitempty"`
	CityId            uint64    `json:"city_id,omitempty"`
	MerchantCreatedAt time.Time `json:"merchant_created_at,omitempty"`
}
