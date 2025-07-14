package models

import "time"

// ProductRecordAttributes represents all the required attributes needed
// to create a product record.
type ProductRecordAttributes struct {
	LastUpdateDate time.Time `json:"last_update_date" validate:"required"`
	PurchasePrice  float64   `json:"purchase_price" validate:"required,gt=0"`
	SalePrice      float64   `json:"sale_price" validate:"required,gtefield=PurchasePrice"`
	ProductID      int       `json:"product_id" validate:"required,gt=0"`
}

// ProductRecord represents a product record, including its
// auto-generated ID plus its attributes.
type ProductRecord struct {
	ID int `json:"id" validate:"required,gt=0"`
	ProductRecordAttributes
}
