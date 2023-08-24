package _interface

import (
	"Ecommerce/internal/pkg/db/models"
)

type ProductRepository interface {
	CreateProduct(product *models.Product) error
	GetByID(id string) (*models.Product, error)
	GetAllProducts() ([]models.Product, error)
	UpdateProduct(id string, product *models.Product) error
	DeleteProduct(id string) error
}
