package models

type Warehouse struct {
	Id int `json:"id"`
	WarehouseAttributes
}
type WarehouseAttributes struct {
	WarehouseCode      string  `json:"warehouse_code"`
	Address            string  `json:"address"`
	Telephone          string  `json:"telephone"`
	MinimunCapacity    int     `json:"minimun_capacity"`
	MinimunTemperature float64 `json:"minimun_temperature"`
}

func (b Warehouse) GetID() int {
	return b.Id
}
