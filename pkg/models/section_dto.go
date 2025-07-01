package models

// CreateSectionRequest es el DTO para la petición POST de una nueva sección.
type CreateSectionRequest struct {
	SectionNumber      string  `json:"section_number" validate:"required"`
	CurrentTemperature float64 `json:"current_temperature" validate:"required"`
	MinimumTemperature float64 `json:"minimum_temperature" validate:"required"`
	CurrentCapacity    int     `json:"current_capacity" validate:"required,gte=0"`
	MinimumCapacity    int     `json:"minimum_capacity" validate:"required,gte=0"`
	MaximumCapacity    int     `json:"maximum_capacity" validate:"required,gte=0"`
	WarehouseID        int     `json:"warehouse_id" validate:"required"`
	ProductTypeID      int     `json:"product_type_id" validate:"required"`
}