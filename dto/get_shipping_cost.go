package dto

type GetShippingCost struct {
	OriginCity      string `json:"origin_city" binding:"required"`
	DestinationCity string `json:"destination_city" binding:"required"`
	Weight          string `json:"weight" binding:"required"`
	Courier         string `json:"courier" binding:"required"`
}
