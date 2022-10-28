package dto

import "kamari.com/model"

type MerchantProfileRes struct {
	Id             uint64            `json:"id"`
	Name           string            `json:"name"`
	Username       string            `json:"username"`
	Description    string            `json:"description"`
	ProfilePicture string            `json:"profile_picture"`
	Address        model.UserAddress `json:"address"`
}
