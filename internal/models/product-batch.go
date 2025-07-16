package models

// ProductBatch represents the core data model for a product batch.
type ProductBatch struct {
    ID                 int     `json:"id"`
	ProductBatchAttibutes
}

// CreateProductBatchRequest defines the structure for creating a new ProductBatch.
// It includes validation tags for go-playground/validator.
type ProductBatchAttibutes struct {
	BatchNumber        int     `json:"batch_number" validate:"required,min=1"`
	CurrentQuantity    int     `json:"current_quantity" validate:"required,gte=0"`
	CurrentTemperature float64 `json:"current_temperature" validate:"required"`
	DueDate            string  `json:"due_date" validate:"required,date_format"`
	InitialQuantity    int     `json:"initial_quantity" validate:"required,gte=0"`
	ManufacturingDate  string  `json:"manufacturing_date" validate:"required,date_format"`
	ManufacturingHour  int     `json:"manufacturing_hour" validate:"required"`
	MinimumTemperature float64 `json:"minimum_temperature" validate:"required"`
	ProductID          int     `json:"product_id" validate:"required,gt=0"`
	SectionID          int     `json:"section_id" validate:"required,gt=0"`
}

// CreateProductBatchRequest defines the structure for creating a new ProductBatch.
// It includes validation tags for go-playground/validator.
type CreateProductBatchRequest struct {
	Data ProductBatchAttibutes `json:"data" validate:"required"`
}