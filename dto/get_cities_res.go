package dto

import "kamari.com/model"

type ProvinceToCity struct {
	Province string       `json:"province"`
	Cities   []model.City `json:"cities"`
}
type GetCitiesResponse struct {
	PerProvinceData []ProvinceToCity `json:"per_province_data"`
}
