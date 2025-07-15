package service

import (
	"context"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
)

type ProductService interface {
	Create(ctx context.Context, product models.ProductAttributes) (models.Product, error)
	GetAll(ctx context.Context) ([]models.Product, error)
	GetByID(ctx context.Context, id int) (models.Product, error)
	Update(ctx context.Context, id int, productAttributes models.ProductPatchRequest) (models.Product, error)
	Delete(ctx context.Context, id int) error
}

type SellerService interface {
	Create(seller models.SellerAttributes) (models.Seller, error)
	GetAll() ([]models.Seller, error)
	GetByID(id int) (models.Seller, error)
	Update(id int, data *models.SellerAttributes) (models.Seller, error)
	Delete(id int) error
}

type LocalityService interface {
	Create(l models.Locality) (models.Locality, error)
	GetByID(id string) (models.Locality, error)
	GetSellerReport(id *string) ([]models.SellerReport, error)
}

// WarehouseService defines warehouse operations.
type WarehouseService interface {
	// GetAll returns all warehouses.
	GetAll() ([]models.Warehouse, error)
	// Create adds a new warehouse.
	Create(warehouseAttributes models.WarehouseAttributes) (models.Warehouse, error)
	// GetByID returns a warehouse by ID.
	GetByID(id int) (models.Warehouse, error)
	// Update modifies a warehouse by ID.
	Update(id int, warehouseAttributes models.WarehouseAttributes) (models.Warehouse, error)
	// Delete removes a warehouse by ID.
	Delete(id int) error
}

type BuyerService interface {
	Create(ctx context.Context, newBuyer models.BuyerAttributes) (models.Buyer, error)
	GetAll(ctx context.Context) ([]models.Buyer, error)
	GetByID(ctx context.Context, id int) (models.Buyer, error)
	Update(ctx context.Context, id int, BuyerData models.BuyerPatchRequest) (models.Buyer, error)
	Delete(ctx context.Context, id int) error
}

type EmployeeService interface {
	Create(Employee models.EmployeeAttributes) (models.Employee, error)
	GetAll() ([]models.Employee, error)
	GetByID(id int) (models.Employee, error)
	Update(id int, employee models.EmployeeAttributes) (models.Employee, error)
	Delete(id int) error
}

type SectionService interface {
	Create(section models.Section) (models.Section, error)
	GetAll() ([]models.Section, error)
	GetByID(id int) (models.Section, error)
	Update(id int, data models.UpdateSectionRequest) (models.Section, error)
	Delete(id int) error
}
