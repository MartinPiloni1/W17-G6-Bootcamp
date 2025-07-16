package service

import (
	"context"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
)

type ProductService interface {
	Create(ctx context.Context, product models.ProductAttributes) (models.Product, error)
	GetAll(ctx context.Context) ([]models.Product, error)
	GetByID(ctx context.Context, id int) (models.Product, error)
	GetRecordsPerProduct(ctx context.Context, id *int) ([]models.ProductRecordCount, error)
	Update(ctx context.Context, id int, productAttributes models.ProductPatchRequest) (models.Product, error)
	Delete(ctx context.Context, id int) error
}

type ProductRecordService interface {
	Create(ctx context.Context, productRecord models.ProductRecordAttributes) (models.ProductRecord, error)
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
	// GetReportByLocalityId retrieves a report of carries by locality ID.
	GetReportByLocalityId(localityId string) ([]models.CarryReport, error)
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
	GetWithPurchaseOrdersCount(ctx context.Context, id *int) ([]models.BuyerWithPurchaseOrdersCount, error)
}

type EmployeeService interface {
	Create(Employee models.EmployeeAttributes) (models.Employee, error)
	GetAll() ([]models.Employee, error)
	GetByID(id int) (models.Employee, error)
	Update(id int, employee models.EmployeeAttributes) (models.Employee, error)
	Delete(id int) error
	ReportInboundOrders(employeeID int) ([]models.EmployeeWithInboundCount, error)
}

type SectionService interface {
	Create(ctx context.Context, section models.Section) (models.Section, error)
	GetAll(ctx context.Context) ([]models.Section, error)
	GetByID(ctx context.Context, id int) (models.Section, error)
	Update(ctx context.Context, id int, data models.UpdateSectionRequest) (models.Section, error)
	Delete(ctx context.Context, id int) error
}

type PurchaseOrderService interface {
	Create(ctx context.Context, newPurchaseOrder models.PurchaseOrderAttributes) (models.PurchaseOrder, error)
}

// CarryService defines operations for managing carries.
type CarryService interface {
	// Create validates and creates a new carry.
	Create(carryAttributes models.CarryAttributes) (models.Carry, error)
}

type InboundOrderService interface {
	Create(attrs models.InboundOrderAttributes) (models.InboundOrder, error)
}
