package _interface

import "Ecommerce/pkg/models"

type AddressRepository interface {
	UpdateAddress(address *models.Address) error
}
