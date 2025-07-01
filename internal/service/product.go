package service

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
)

type ProductService struct {
	repo repository.ProductRepositoryInterface
}

func NewProductService(repo repository.ProductRepositoryInterface) ProductServiceInterface {
	return &ProductService{repo: repo}
}

// Create implements ProductServiceInterface.
func (p *ProductService) Create(product models.ProductAtributtes) (models.Product, error) {
	return p.repo.Create(product)
}

// GetAll implements ProductServiceInterface.
func (p *ProductService) GetAll() (map[int]models.Product, error) {
	return p.repo.GetAll()
}

// GetByID implements ProductServiceInterface.
func (p *ProductService) GetByID(id int) (models.Product, error) {
	return p.repo.GetByID(id)
}

// Update implements ProductServiceInterface.
func (p *ProductService) Update(id int, productAtributtes models.ProductAtributtes) (models.Product, error) {
	return p.repo.Update(id, productAtributtes)
}

// Delete implements ProductServiceInterface.
func (p *ProductService) Delete(id int) error {
	return p.repo.Delete(id)
}
