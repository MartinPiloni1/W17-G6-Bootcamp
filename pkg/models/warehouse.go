package models

type Warehouse struct {
	Id                 int    `json:"id"`
	WarehouseCode      string  `json:"warehouse_code"`
	Address            string  `json:"address"`
	Telephone          string  `json:"telephone"`
	MinimunCapacity    int      `json:"minimun_capacity"`
	MinimunTemperature float64  `json:"minimun_temperature"`
}
ß