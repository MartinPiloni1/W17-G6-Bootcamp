package models

// ProductAttributes represents all the required attributes needed
// to create a product.
type ProductAttributes struct {
	Description                    string  `json:"description" validate:"required,min=5,max=500"`
	ExpirationRate                 int     `json:"expiration_rate" validate:"required,gte=0"`
	FreezingRate                   int     `json:"freezing_rate" validate:"required,gte=0"`
	Height                         float64 `json:"height" validate:"required,gt=0"`
	Length                         float64 `json:"length" validate:"required,gt=0"`
	Width                          float64 `json:"width" validate:"required,gt=0"`
	NetWeight                      float64 `json:"netweight" validate:"required,gt=0"`
	ProductCode                    string  `json:"product_code" validate:"required,alphanum,min=1,max=50"`
	RecommendedFreezingTemperature float64 `json:"recommended_freezing_temperature" validate:"required,gte=-80,lte=25"`
	ProductTypeID                  int     `json:"product_type_id" validate:"required,gt=0"`
	SellerID                       int     `json:"seller_id,omitempty" validate:"omitempty,gt=0"`
}

// ProductPatchRequest holds optional fields for partial updates.
// Fields are pointers so that we can diferentiate empty fields and fields
// with default values.
type ProductPatchRequest struct {
	Description                    *string  `json:"description,omitempty" validate:"omitempty,min=5,max=500"`
	ExpirationRate                 *int     `json:"expiration_rate,omitempty" validate:"omitempty,gte=0"`
	FreezingRate                   *int     `json:"freezing_rate,omitempty" validate:"omitempty,gte=0"`
	Height                         *float64 `json:"height,omitempty" validate:"omitempty,gt=0"`
	Length                         *float64 `json:"length,omitempty" validate:"omitempty,gt=0"`
	Width                          *float64 `json:"width,omitempty" validate:"omitempty,gt=0"`
	NetWeight                      *float64 `json:"netweight,omitempty" validate:"omitempty,gt=0"`
	ProductCode                    *string  `json:"product_code,omitempty" validate:"omitempty,alphanum,min=1,max=50"`
	RecommendedFreezingTemperature *float64 `json:"recommended_freezing_temperature,omitempty" validate:"omitempty,gte=-80,lte=25"`
	ProductTypeID                  *int     `json:"product_type_id,omitempty" validate:"omitempty,gt=0"`
	SellerID                       *int     `json:"seller_id,omitempty" validate:"omitempty,gt=0"`
}

// Product represents a product entry, including its automatically assigned ID
// plus all attributes.
type Product struct {
	ID int `json:"id"`
	ProductAttributes
}
