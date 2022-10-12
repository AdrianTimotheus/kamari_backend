package dto

type AddNewAddressRequest struct {
	Title         string `json:"title" binding:"required"`
	RecipientName string `json:"recipient_name" binding:"required"`
	Phone         string `json:"phone" binding:"required"`
	CityId        uint64 `json:"city_id" binding:"required"`
	DetailAddress string `json:"detail_address" binding:"required"`
}
