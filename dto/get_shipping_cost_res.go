package dto

type CourierServiceRes struct {
	Code  string               `json:"code"`
	Name  string               `json:"name"`
	Costs []CourierServiceCost `json:"costs"`
}

type CourierServiceCost struct {
	Service     string              `json:"service"`
	Description string              `json:"description"`
	Cost        []ServiceCostDetail `json:"cost"`
}

type ServiceCostDetail struct {
	Value int    `json:"value"`
	Etd   string `json:"etd"`
	Note  string `json:"note"`
}
