package models

import (
	"time"
)

type PurchaseOrder struct {
	Id int `json:"id"`
	PurchaseOrderAttributes
}

type PurchaseOrderAttributes struct {
	OrderNumber     string    `json:"order_number" validate:"required,max=255,aphanum"`
	OrderDate       time.Time `json:"order_date" validate:"required"`
	TrackingCode    string    `json:"tracking_code" validate:"required,max=255,alphanum"`
	BuyerId         int       `json:"buyer_id" validate:"required,gt=0"`
	ProductRecordId int       `json:"product_record_id" validate:"required,gt=0"`
}

type BuyerWithPurchaseOrdersCount struct {
	Buyer
	PurchaseOrdersCount int `json:"purchase_orders_count"`
}
