package _interface

import (
	"Ecommerce/internal/pkg/db/models"
)

type AddressRepository interface {
	UpdateAddress(address *models.Address) error
}
