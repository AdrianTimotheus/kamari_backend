package dto

import "kamari.com/model"

type CategoryResponse struct {
	Id            uint64              `json:"id"`
	Name          string              `json:"name"`
	PicUrl        string              `json:"pic_url"`
	Subcategories []model.Subcategory `json:"subcategories"`
}
