package models

type Buyer struct {
	Id int `json:"id"`
	BuyerAttributes
}

type BuyerAttributes struct {
	CardNumberId int    `json:"card_number_id" validate:"required,gte=10000000,lte=99999999"` // like a dni: 8 digits
	FirstName    string `json:"first_name" validate:"required,min=1,max=100,alpha"`           // only letters
	LastName     string `json:"last_name" validate:"required,min=1,max=100,alpha"`            // only letters
}

type BuyerPatchRequest struct {
	CardNumberId *int    `json:"card_number_id" validate:"omitempty,gte=10000000,lte=99999999"` // like a dni: 8 digits
	FirstName    *string `json:"first_name" validate:"omitempty,min=1,max=100,alpha"`           // only letters
	LastName     *string `json:"last_name" validate:"omitempty,min=1,max=100,alpha"`            // only letters
}

type BuyerWithPurchaseOrdersCount struct {
	Buyer
	PurchaseOrdersCount int `json:"purchase_orders_count"`
}
