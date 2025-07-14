package models

// Warehouse represents a warehouse entity with its ID and attributes.
type Warehouse struct {
	Id int `json:"id"`
	WarehouseAttributes
}

// WarehouseAttributes holds the details of a warehouse.
type WarehouseAttributes struct {
	WarehouseCode      string  `json:"warehouse_code"`
	Address            string  `json:"address"`
	Telephone          string  `json:"telephone"`
	MinimunCapacity    int     `json:"minimun_capacity"`
	MinimunTemperature float64 `json:"minimun_temperature"`
}
