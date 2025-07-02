package models

type Employee struct {
	Id int `json:"id"`
	EmployeeAttributes
}

type EmployeeAttributes struct {
	CardNumberID string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseID  int    `json:"warehouse_id"`
}

func (e Employee) GetID() int {
	return e.Id
}
