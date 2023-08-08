package models

import (
	"gorm.io/gorm"
	"time"
)

//нврн надо менять

type User struct {
	gorm.Model
	Username  string            `json:"username"`
	Email     string            `gorm:"uniqueIndex" json:"email" binding:"required" validate:"required,email"`
	Password  string            `json:"password" binding:"required" validate:"required,min=8,max=64"`
	UserType  string            `json:"userType"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
	Address   Address           `gorm:"foreignKey:ZipCode" json:"address"`
	Orders    []Order           `gorm:"foreignKey:OrderId"json:"orders"`
	UserCart  []ProductsToOrder `gorm:"foreignKey:ProductId"json:"userCart"`
}

type Product struct {
	gorm.Model
	Name              string   `gorm:"not null,uniqueIndex" json:"name,omitempty" validate:"required"`
	Price             float64  `gorm:"not null"  json:"price,omitempty" validate:"required"`
	Description       string   `gorm:"not null"  json:"description,omitempty" validate:"required"`
	AvailableQuantity int      `gorm:"not null"  json:"availableQuantity" validate:"required"`
	Category          string   `gorm:"not null" json:"category,omitempty" validate:"required"`
	Images            []string `json:"images,omitempty"`
}

type Order struct {
	gorm.Model
	OrderCart  []ProductsToOrder `gorm:"foreignKey:ProductId" json:"orderCart,omitempty"`
	TotalPrice float64           `json:"totalPrice,omitempty"`
}

type ProductsToOrder struct {
	gorm.Model
	Name        string  `gorm:"not null" json:"name,omitempty"`
	Price       float64 `gorm:"not null" json:"price,omitempty"`
	BuyQuantity int     `json:"buyQuantity"`
}

type Address struct {
	gorm.Model
	ZipCode     string `gorm:"primaryKey" json:"zipCode,omitempty" validate:"required"`
	City        string `gorm:"not null" json:"city,omitempty" validate:"required"`
	State       string `gorm:"not null" json:"state,omitempty" validate:"required"`
	Country     string `gorm:"not null" json:"country,omitempty" validate:"required"`
	Street      string `gorm:"not null" json:"street,omitempty" validate:"required"`
	HouseNumber string `gorm:"not null" json:"houseNumber,omitempty" validate:"required"`
}
