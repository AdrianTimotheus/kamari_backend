package dto

type ChangePinWalletRequest struct {
	Password string `json:"password" binding:"required"`
}
