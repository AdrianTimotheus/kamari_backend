package dto

type MerchantProfileReq struct {
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description" binding:"required"`
	ProfilePicture string `json:"profile_picture" binding:"required"`
}
