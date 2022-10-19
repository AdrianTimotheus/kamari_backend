package dto

type AllCouriersResponse struct {
	Couriers []SingleCourierResponse `json:"couriers"`
}

type SingleCourierResponse struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}
