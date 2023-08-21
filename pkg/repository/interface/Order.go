package _interface

import "Ecommerce/pkg/models"

type OrderRepository interface {
	CreateOrder(order *models.Order) error
	DeleteProductToOrder(PO *models.ProductsToOrder) error
	UpdateProductToOrder(id string, PO *models.ProductsToOrder) error
}
