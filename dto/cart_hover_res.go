package dto

type CartHoverRes struct {
	RemainingItem uint64          `json:"remaining_item"`
	Items         []CartHoverItem `json:"items"`
}

type CartHoverItem struct {
	Name   string `json:"name"`
	Price  int    `json:"price"`
	PicUrl string `json:"pic_url"`
}
