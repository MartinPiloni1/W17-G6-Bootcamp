package service

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
)

type ProductService struct {
	repo repository.ProductRepositoryInterface
}

func NewProductService(repo repository.ProductRepositoryInterface) ProductServiceInterface {
	return &ProductService{repo: repo}
}

// Create implements ProductServiceInterface.
func (p *ProductService) Create(Product models.Product) (*models.Product, error) {
	panic("unimplemented")
}

// GetAll implements ProductServiceInterface.
func (p *ProductService) GetAll() (map[int]models.Product, error) {
	return p.repo.GetAll()
}

// GetByID implements ProductServiceInterface.
func (p *ProductService) GetByID(id int) (models.Product, error) {
	return p.repo.GetByID(id)
}

// Delete implements ProductServiceInterface.
func (p *ProductService) Delete(id int) error {
	panic("unimplemented")
}

// Update implements ProductServiceInterface.
func (p *ProductService) Update(id int, data models.Product) (models.Product, error) {
	panic("unimplemented")
}
