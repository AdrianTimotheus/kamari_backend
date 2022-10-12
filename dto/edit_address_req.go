package dto

type EditAddressReq struct {
	Id            uint64 `json:"id" binding:"required"`
	Title         string `json:"title"`
	RecipientName string `json:"recipient_name"`
	Phone         string `json:"phone"`
	CityId        uint64 `json:"city_id"`
	DetailAddress string `json:"detail_address"`
}
