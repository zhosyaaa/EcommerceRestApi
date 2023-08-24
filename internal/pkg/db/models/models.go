package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string            `json:"username" `
	Email    string            `json:"email" `
	Password string            `json:"password"`
	UserType string            `json:"userType"`
	Address  Address           `json:"address" gorm:"foreignKey:ID"`
	Orders   []Order           `json:"orders" gorm:"foreignKey:ID"`
	UserCart []ProductsToOrder `json:"userCart" gorm:"foreignKey:ID"`
}

type Product struct {
	gorm.Model
	Name              string  `json:"name,omitempty"`
	Price             float64 `json:"price,omitempty"`
	Description       string  `json:"description,omitempty"`
	AvailableQuantity int     `json:"availableQuantity"`
	Category          string  `json:"category,omitempty"`
	Image             string  `json:"image,omitempty"`
}

type Order struct {
	gorm.Model
	OrderCart  []ProductsToOrder `json:"orderCart,omitempty" gorm:"foreignKey:ID"`
	TotalPrice float64           `json:"totalPrice,omitempty"`
	AddressID  uint              `json:"addressID" gorm:"foreignKey:ID"`
}

type ProductsToOrder struct {
	gorm.Model
	Name        string  `json:"name,omitempty"`
	Price       float64 `json:"price,omitempty"`
	BuyQuantity int     `json:"buyQuantity"`
}

type Address struct {
	gorm.Model
	ZipCode     string `json:"zipCode,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	Country     string `json:"country,omitempty"`
	Street      string `json:"street,omitempty"`
	HouseNumber string `json:"houseNumber,omitempty"`
}

type AddressInputCred struct {
	ZipCode     string `json:"zipCode"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	Street      string `json:"street"`
	HouseNumber string `json:"houseNumber"`
}
