package service

import (
	interfaces "Ecommerce/internal/app/service/interface"
	"Ecommerce/internal/pkg/db/models"
	"gorm.io/gorm"
)

type addressDatabase struct {
	DB *gorm.DB
}

func (a addressDatabase) UpdateAddress(address *models.Address) error {
	if err := a.DB.Save(&address).Error; err != nil {
		return err
	}
	return nil
}

func NewAddressDatabase(DB *gorm.DB) interfaces.AddressRepository {
	return &addressDatabase{DB: DB}
}
