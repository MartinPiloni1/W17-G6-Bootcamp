package models

type Buyer struct {
	Id int `json:"id"`
	BuyerAttributes
}

type BuyerAttributes struct {
	CardNumberId int    `json:"card_number_id" validate:"required,gt=0"`
	FirstName    string `json:"first_name" validate:"required,min=1"`
	LastName     string `json:"last_name" validate:"required,min=1"`
}

func (b Buyer) GetID() int {
	return b.Id
}

type BuyerPatchRequest struct {
	CardNumberId *int    `json:"card_number_id" validate:"omitempty,gt=0"`
	FirstName    *string `json:"first_name" validate:"omitempty,min=1"`
	LastName     *string `json:"last_name" validate:"omitempty,min=1"`
}
