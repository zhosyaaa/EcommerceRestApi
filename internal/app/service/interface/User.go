package _interface

import (
	"Ecommerce/internal/pkg/db/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	UpdateUser(id string, user *models.User) error
	GetAllUsers(users []models.User) ([]models.User, error)
	DeleteUser(user *models.User) error
	DeleteAllUsers(users []models.User) error
}
