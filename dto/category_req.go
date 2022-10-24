package dto

type CategoryRequest struct {
	Name   string `json:"name" binding:"required"`
	PicUrl string `json:"picUrl" binding:"required"`
}
