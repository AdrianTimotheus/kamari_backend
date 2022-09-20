package handler

import "kamari.com/service"

type Handler struct {
	provinceService service.ProvinceService
}

type HConfig struct {
	ProvinceService service.ProvinceService
}

func New(c HConfig) Handler {
	return Handler{
		provinceService: c.ProvinceService,
	}
}
