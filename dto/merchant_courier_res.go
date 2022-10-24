package dto

type MerchantCourierRes struct {
	MerchantId    int                        `json:"merchantId"`
	Weight        int                        `json:"weight"`
	CourierDetail []MerchantCourierDetailRes `json:"availableCourier"`
}

type MerchantCourierDetailRes struct {
	Courier string  `json:"courier"`
	Price   float64 `json:"price"`
	Etd     string  `json:"etd"`
}
