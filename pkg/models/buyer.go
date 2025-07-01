package models

type Buyer struct {
	Id int `json:"id"`
	BuyerAttributes
}

type BuyerAttributes struct {
	CardNumberId int    `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

func (b Buyer) GetID() int {
	return b.Id
}

type BuyerPatchRequest struct {
	CardNumberId *int    `json:"card_number_id"`
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
}
