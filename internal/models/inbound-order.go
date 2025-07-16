package models

import "time"

// InboundOrder represents an inbound order in the system.
type InboundOrder struct {
	ID int `json:"id"`
	InboundOrderAttributes
}

type InboundOrderAttributes struct {
	OrderNumber    string    `json:"order_number" validate:"required"`
	OrderDate      time.Time `json:"order_date"`
	EmployeeID     int       `json:"employee_id" validate:"required,gt=0"`
	WarehouseID    int       `json:"warehouse_id" validate:"required,gt=0"`
	ProductBatchID int       `json:"product_batch_id" validate:"required,gt=0"`
}

type EmployeeWithInboundCount struct {
	Id                 int    `json:"id"`
	CardNumberID       string `json:"card_number_id"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	WarehouseID        int    `json:"warehouse_id"`
	InboundOrdersCount int    `json:"inbound_orders_count"`
}
