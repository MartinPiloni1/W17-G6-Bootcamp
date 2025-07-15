package models

type SellerAttributes struct {
	CID         int    `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
	LocalityID  string `json:"locality_id"`
}
type Seller struct {
	ID int `json:"id"`
	SellerAttributes
}

func (s Seller) GetID() int {
	return s.ID
}
