package models

type Section struct {
	ID                 int            `json:"id"`
	SectionNumber      string         `json:"section_number"`
	CurrentTemperature float64        `json:"current_temperature"`
	MinimumTemperature float64        `json:"minimum_temperature"`
	CurrentCapacity    int            `json:"current_capacity"`
	MinimumCapacity    int            `json:"minimum_capacity"`
	MaximumCapacity    int            `json:"maximum_capacity"`
	WarehouseID        int            `json:"warehouse_id"`
	ProductTypeID      int            `json:"product_type_id"`
	// ProductBatches     []	 `json:"product_batches"`
}

func (s Section) GetID() int {
	return s.ID
}

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