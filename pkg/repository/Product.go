package repository

import (
	"Ecommerce/pkg/models"
	interfaces "Ecommerce/pkg/repository/interface"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type productDatabase struct {
	DB *gorm.DB
}

func NewProductRepository(DB *gorm.DB) interfaces.ProductRepository {
	return &productDatabase{DB}
}

func (p productDatabase) CreateProduct(product *models.Product) error {
	if err := p.DB.Create(&product).Error; err != nil {
		return err
	}
	return nil
}

func (p productDatabase) GetByID(id string) (*models.Product, error) {
	var product models.Product
	if err := p.DB.Where("ID=?", id).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (p productDatabase) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	err := p.DB.Find(&products).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching products: %w", err)
	}
	return products, nil
}

func (p productDatabase) UpdateProduct(id string, product *models.Product) error {
	var existingProduct models.Product
	if err := p.DB.First(&existingProduct, id).Error; err != nil {
		return err
	}
	existingProduct.Name = product.Name
	existingProduct.Description = product.Description
	existingProduct.Price = product.Price
	existingProduct.AvailableQuantity = product.AvailableQuantity
	existingProduct.Category = product.Category
	existingProduct.Image = product.Image

	if err := p.DB.Save(&existingProduct).Error; err != nil {
		return err
	}

	return nil
}

func (p productDatabase) DeleteProduct(id string) error {
	var product models.Product
	if err := p.DB.First(&product, id).Error; err != nil {
		return err
	}

	// Delete the product from the database
	if err := p.DB.Delete(&product).Error; err != nil {
		return err
	}

	return nil
}
