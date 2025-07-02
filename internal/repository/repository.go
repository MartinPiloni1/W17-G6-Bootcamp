package repository

import "github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"

type SellerRepository interface {
	Create(seller models.SellerAttributes) (models.Seller, error)
	GetAll() (map[int]models.Seller, error)
	GetByID(id int) (models.Seller, error)
	Update(id int, data *models.SellerAttributes) (models.Seller, error)
	Delete(id int) error
}
type WarehouseRepository interface {
	GetAll() (map[int]models.Warehouse, error)
	Create(warehouseAtribbutes models.WarehouseAttributes) (models.Warehouse, error)
	GetByID(id int) (models.Warehouse, error)
	Update(id int, warehouseAttributes models.WarehouseAttributes) (models.Warehouse, error)
	Delete(id int) error
}

type BuyerRepository interface {
	Create(newBuyer models.BuyerAttributes) (models.Buyer, error)
	GetAll() (map[int]models.Buyer, error)
	GetByID(id int) (models.Buyer, error)
	Update(id int, updatedBuyer models.Buyer) (models.Buyer, error)
	Delete(id int) error
	CardNumberIdAlreadyExist(newCardNumberId int) (bool, error)
}

type EmployeeRepository interface {
	Create(Employee models.Employee) (models.Employee, error)
	GetAll() (map[int]models.Employee, error)
	GetByID(id int) (models.Employee, error)
	Update(id int, data models.Employee) (models.Employee, error)
	Delete(id int) error
}