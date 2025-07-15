package models

// Employee represents an employee entity with its unique ID and all core attributes.
type Employee struct {
	Id int `json:"id"`
	EmployeeAttributes
}

// EmployeeAttributes contains the fields necessary to create or update an employee.
type EmployeeAttributes struct {
	CardNumberID string `json:"card_number_id" validate:"required,len=8,numeric"`
	FirstName    string `json:"first_name" validate:"required,min=1,max=100,alpha"`
	LastName     string `json:"last_name" validate:"required,min=1,max=100,alpha"`
	WarehouseID  int    `json:"warehouse_id" validate:"required,gt=0"`
}
