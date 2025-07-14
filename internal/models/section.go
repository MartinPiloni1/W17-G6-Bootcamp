package models

// In internal/models/section.go

// Section represents the core data model for a section.
type Section struct {
    ID                 int     `json:"id"`
    SectionNumber      string  `json:"section_number"`
    CurrentTemperature float64 `json:"current_temperature"`
    MinimumTemperature float64 `json:"minimum_temperature"`
    CurrentCapacity    int     `json:"current_capacity"`
    MinimumCapacity    int     `json:"minimum_capacity"`
    MaximumCapacity    int     `json:"maximum_capacity"`
    WarehouseID        int     `json:"warehouse_id"`
    ProductTypeID      int     `json:"product_type_id"`
}

// CreateSectionRequest defines the structure for creating a new Section.
// It includes validation tags for go-playground/validator.
type CreateSectionRequest struct {
    SectionNumber      string  `json:"section_number" validate:"required,min=1"`
    CurrentTemperature float64 `json:"current_temperature" validate:"required"`
    MinimumTemperature float64 `json:"minimum_temperature" validate:"required"`
    CurrentCapacity    int     `json:"current_capacity" validate:"required,gte=0"`
    MinimumCapacity    int     `json:"minimum_capacity" validate:"required,gte=0"`
    MaximumCapacity    int     `json:"maximum_capacity" validate:"required,gte=0"`
    WarehouseID        int     `json:"warehouse_id" validate:"required,gt=0"`
    ProductTypeID      int     `json:"product_type_id" validate:"required,gt=0"`
}

// UpdateSectionRequest defines the structure for updating an existing Section.
// Fields are pointers to allow for partial updates (omitempty).
// It includes validation tags for go-playground/validator.
type UpdateSectionRequest struct {
    SectionNumber      *string  `json:"section_number,omitempty" validate:"omitempty,min=1"`
    CurrentTemperature *float64 `json:"current_temperature,omitempty"`
    MinimumTemperature *float64 `json:"minimum_temperature,omitempty"`
    CurrentCapacity    *int     `json:"current_capacity,omitempty" validate:"omitempty,gte=0"`
    MinimumCapacity    *int     `json:"minimum_capacity,omitempty" validate:"omitempty,gte=0"`
    MaximumCapacity    *int     `json:"maximum_capacity,omitempty" validate:"omitempty,gte=0"`
    WarehouseID        *int     `json:"warehouse_id,omitempty" validate:"omitempty,gt=0"`
    ProductTypeID      *int     `json:"product_type_id,omitempty" validate:"omitempty,gt=0"`
}