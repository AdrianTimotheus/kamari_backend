package dto

type MerchantDetailResponse struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	ProfilePicture string `json:"profile_picture"`
	FollowerCount  uint64 `json:"follower_count"`
	ProductCount   uint64 `json:"product_count"`
	ShopLocation   string `json:"shop_location"`
}

//seller detail
//- show merchant name
//- show merchant description
//- get_shop_detail&username=unique username
//- shop location
//- show product count
//- show follower count
