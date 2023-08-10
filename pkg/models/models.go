package models

import (
	"time"
)

type User struct {
	ID        int               `json:"ID" gorm:"primarykey"`
	Username  string            `json:"username" `
	Email     string            `json:"email" `
	Password  string            `json:"password"`
	UserType  string            `json:"userType"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
	Address   Address           `json:"address" gorm:"foreignKey:ID"`
	Orders    []Order           `json:"orders" gorm:"foreignKey:Id"`
	UserCart  []ProductsToOrder `json:"userCart" gorm:"foreignKey:ProductId"`
}

type Product struct {
	ID                int       `json:"_id,omitempty" gorm:"primarykey"`
	Name              string    `json:"name,omitempty"`
	Price             float64   `json:"price,omitempty"`
	Description       string    `json:"description,omitempty"`
	AvailableQuantity int       `json:"availableQuantity"`
	Category          string    `json:"category,omitempty"`
	Images            []string  `json:"images,omitempty"`
	CreatedAt         time.Time `json:"createdAt,omitempty"`
	UpdatedAt         time.Time `json:"updatedAt,omitempty"`
}

type Order struct {
	Id         int               `json:"_id,omitempty" gorm:"primarykey"`
	OrderCart  []ProductsToOrder `json:"orderCart,omitempty" gorm:"foreignKey:ProductId"`
	TotalPrice float64           `json:"totalPrice,omitempty"`
	CreatedAt  time.Time         `json:"createdAt,omitempty"`
	AddressID  int               `json:"addressID" gorm:"foreignKey:AddressID"`
}

type ProductsToOrder struct {
	ProductId   int       `json:"productId,omitempty" gorm:"primarykey"`
	Name        string    `json:"name,omitempty"`
	Price       float64   `json:"price,omitempty"`
	BuyQuantity int       `json:"buyQuantity"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

type Address struct {
	ID          int    `json:"id" gorm:"primarykey"`
	ZipCode     string `json:"zipCode,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	Country     string `json:"country,omitempty"`
	Street      string `json:"street,omitempty"`
	HouseNumber string `json:"houseNumber,omitempty"`
}
