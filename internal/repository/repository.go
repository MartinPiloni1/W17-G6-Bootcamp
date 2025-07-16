package repository

import (
	"context"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
)

type ProductRepository interface {
	Create(ctx context.Context, product models.ProductAttributes) (models.Product, error)
	GetAll(ctx context.Context) ([]models.Product, error)
	GetByID(ctx context.Context, id int) (models.Product, error)
	GetRecordsPerProduct(ctx context.Context, id *int) ([]models.ProductRecordCount, error)
	Update(ctx context.Context, id int, product models.Product) (models.Product, error)
	Delete(ctx context.Context, id int) error
}

type ProductRecordRepository interface {
	Create(ctx context.Context, productRecord models.ProductRecordAttributes) (models.ProductRecord, error)
}

type SellerRepository interface {
	Create(seller models.SellerAttributes) (models.Seller, error)
	GetAll() ([]models.Seller, error)
	GetByID(id int) (models.Seller, error)
	Update(id int, data *models.SellerAttributes) (models.Seller, error)
	Delete(id int) error
}

type LocalityRepository interface {
	Create(locality models.Locality) (models.Locality, error)
	GetByID(id string) (models.Locality, error)
	GetSellerReport(id *string) ([]models.SellerReport, error)
	// GetReportByLocalityId retrieves a report of carries by locality ID.
	GetReportByLocalityId(localityId string) ([]models.CarryReport, error)
}

// WarehouseRepository provides methods for warehouse data access.
type WarehouseRepository interface {
	// GetAll returns all warehouses.
	GetAll() ([]models.Warehouse, error)
	// Create adds a new warehouse.
	Create(warehouseAtribbutes models.WarehouseAttributes) (models.Warehouse, error)
	// GetByID returns a warehouse by its ID.
	GetByID(id int) (models.Warehouse, error)
	// Update modifies a warehouse by its ID.
	Update(id int, warehouseAttributes models.WarehouseAttributes) (models.Warehouse, error)
	// Delete removes a warehouse by its ID.
	Delete(id int) error
}

type BuyerRepository interface {
	Create(ctx context.Context, newBuyer models.BuyerAttributes) (models.Buyer, error)
	GetAll(ctx context.Context) ([]models.Buyer, error)
	GetByID(ctx context.Context, id int) (models.Buyer, error)
	Update(ctx context.Context, id int, updatedBuyer models.Buyer) (models.Buyer, error)
	Delete(ctx context.Context, id int) error
}

type EmployeeRepository interface {
	Create(Employee models.Employee) (models.Employee, error)
	GetAll() ([]models.Employee, error)
	GetByID(id int) (models.Employee, error)
	Update(id int, data models.Employee) (models.Employee, error)
	Delete(id int) error
}

type SectionRepository interface {
	Create(ctx context.Context, section models.Section) (models.Section, error)
	GetAll(ctx context.Context) ([]models.Section, error)
	GetByID(ctx context.Context, id int) (models.Section, error)
	Update(ctx context.Context, id int, data models.Section) (models.Section, error)
	Delete(ctx context.Context, id int) error
}

// CarryRepository provides methods for carry data access.
type CarryRepository interface {
	// Create creates a new carry.
	Create(carryAttributes models.CarryAttributes) (models.Carry, error)
}
