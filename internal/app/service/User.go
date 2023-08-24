package service

import (
	interfaces "Ecommerce/internal/app/service/interface"
	"Ecommerce/internal/pkg/db/models"
	"errors"
	"gorm.io/gorm"
)

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userDatabase{DB}
}

func (u *userDatabase) DeleteUser(user *models.User) error {
	if err := u.DB.Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

func (u *userDatabase) DeleteAllUsers(users []models.User) error {
	if err := u.DB.Delete(&users).Error; err != nil {
		return err
	}
	return nil
}

func (u *userDatabase) GetAllUsers(users []models.User) ([]models.User, error) {
	if err := u.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userDatabase) UpdateUser(id string, user *models.User) error {
	if err := u.DB.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (u *userDatabase) CreateUser(user *models.User) error {
	if err := u.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (u *userDatabase) GetByID(id string) (*models.User, error) {
	var user models.User
	if err := u.DB.Where("ID=?", id).Preload("Address").Preload("UserCart").Preload("Orders").First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}
func (u *userDatabase) GetByEmail(email string) (*models.User, error) {
	var existingUser models.User
	result := u.DB.Where("email = ?", email).Preload("Address").Preload("UserCart").Preload("Orders").First(&existingUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return &existingUser, nil
}
