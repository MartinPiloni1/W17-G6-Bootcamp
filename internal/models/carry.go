package models

// Carry represents a carry entity with its unique ID and attributes.
type Carry struct {
	Id int `json:"id"`
	CarryAttributes
}

// CarryAttributes holds the details of a carry, such as company name, address, telephone, and locality ID.
type CarryAttributes struct {
	Cid         string `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
	LocalityId  string `json:"locality_id"`
}

// CarryReport is used for reporting purposes, containing locality information and the count of carries.
type CarryReport struct {
	LocalityId   string `json:"locality_id"`
	LocalityName string `json:"locality_name"`
	CarriesCount int    `json:"carries_count"`
}
