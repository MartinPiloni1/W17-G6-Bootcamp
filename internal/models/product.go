package models

type Product struct {
	ID int `json:"id"`
	ProductAttributes
}

type ProductAttributes struct {
	Description                    string  `json:"description"`
	ExpirationRate                 int     `json:"expiration_rate"`
	FreezingRate                   int     `json:"freezing_rate"`
	Height                         float64 `json:"height"`
	Length                         float64 `json:"length"`
	Width                          float64 `json:"width"`
	NetWeight                      float64 `json:"netweight"`
	ProductCode                    string  `json:"product_code"`
	RecommendedFreezingTemperature float64 `json:"recommended_freezing_temperature"`
	ProductTypeID                  int     `json:"product_type_id"`
	SellerID                       int     `json:"seller_id"`
}

func (p Product) GetID() int {
	return p.ID
}
