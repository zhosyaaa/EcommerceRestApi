package repository

import (
	"Ecommerce/pkg/models"
	interfaces "Ecommerce/pkg/repository/interface"
	"gorm.io/gorm"
)

type orderDatabase struct {
	DB *gorm.DB
}

func NewOrderRepository(DB *gorm.DB) interfaces.OrderRepository {
	return &orderDatabase{DB: DB}
}

func (o orderDatabase) CreateOrder(order *models.Order) error {
	if err := o.DB.Create(&order).Error; err != nil {
		return err
	}
	return nil
}

func (o orderDatabase) UpdateProductToOrder(id string, PO *models.ProductsToOrder) error {
	var existingPO models.ProductsToOrder
	if err := o.DB.First(&existingPO, id).Error; err != nil {
		return err
	}

	existingPO.Name = PO.Name
	existingPO.Price = PO.Price
	existingPO.BuyQuantity = PO.BuyQuantity

	if err := o.DB.Save(&existingPO).Error; err != nil {
		return err
	}
	return nil
}

func (o orderDatabase) DeleteProductToOrder(PO *models.ProductsToOrder) error {
	var po models.ProductsToOrder
	if err := o.DB.First(&po, PO.ID).Error; err != nil {
		return err
	}

	if err := o.DB.Delete(&po).Error; err != nil {
		return err
	}

	return nil
}
