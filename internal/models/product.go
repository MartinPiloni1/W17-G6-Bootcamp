package models

type ProductAttributes struct {
	Description                    string  `json:"description" validate:"required,min=5,max=500"`
	ExpirationRate                 int     `json:"expiration_rate" validate:"required,gte=0"`
	FreezingRate                   int     `json:"freezing_rate" validate:"required,gte=0"`
	Height                         float64 `json:"height" validate:"required,gte=0"`
	Length                         float64 `json:"length" validate:"required,gte=0"`
	Width                          float64 `json:"width" validate:"required,gte=0"`
	NetWeight                      float64 `json:"netweight" validate:"required,gt=0"`
	ProductCode                    string  `json:"product_code" validate:"required,alphanum,min=1,max=50"`
	RecommendedFreezingTemperature float64 `json:"recommended_freezing_temperature" validate:"required"`
	ProductTypeID                  int     `json:"product_type_id" validate:"required,gt=0"`
	SellerID                       int     `json:"seller_id,omitempty" validate:"omitempty,gt=0"`
}

// Usamos punteros para poder diferenciar entre un campo no enviado (nil)
// y un campo enviado con su valor cero (ej: 0, "").
// Las etiquetas de validación usan 'omitempty' para que solo se apliquen
// si el campo está presente en la petición.
type ProductPatchRequest struct {
	Description                    *string  `json:"description,omitempty" validate:"omitempty,min=5,max=500"`
	ExpirationRate                 *int     `json:"expiration_rate,omitempty" validate:"omitempty,gte=0"`
	FreezingRate                   *int     `json:"freezing_rate,omitempty" validate:"omitempty,gte=0"`
	Height                         *float64 `json:"height,omitempty" validate:"omitempty,gte=0"`
	Length                         *float64 `json:"length,omitempty" validate:"omitempty,gte=0"`
	Width                          *float64 `json:"width,omitempty" validate:"omitempty,gte=0"`
	NetWeight                      *float64 `json:"netweight,omitempty" validate:"omitempty,gt=0"`
	ProductCode                    *string  `json:"product_code,omitempty" validate:"omitempty,alphanum,min=1,max=50"`
	RecommendedFreezingTemperature *float64 `json:"recommended_freezing_temperature,omitempty" validate:"omitempty"`
	ProductTypeID                  *int     `json:"product_type_id,omitempty" validate:"omitempty,gt=0"`
	SellerID                       *int     `json:"seller_id,omitempty" validate:"omitempty,gt=0"`
}

type Product struct {
	ID int `json:"id"`
	ProductAttributes
}

func (p Product) GetID() int {
	return p.ID
}
