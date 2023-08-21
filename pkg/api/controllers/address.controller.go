package controllers

import (
	"Ecommerce/pkg/models"
	interfaces "Ecommerce/pkg/repository/interface"
	"github.com/gin-gonic/gin"
)

type AddressController struct {
	addressService interfaces.AddressRepository
	userService    interfaces.UserRepository
}

func NewAddressController(addressService interfaces.AddressRepository, userService interfaces.UserRepository) *AddressController {
	return &AddressController{addressService: addressService, userService: userService}
}

func (s *AddressController) UpdateAddress(c *gin.Context) {
	userId := c.Param("id")
	var addressInput models.AddressInputCred
	if err := c.ShouldBindJSON(&addressInput); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid address data",
			"data":    err.Error(),
		})
		return
	}
	user, err := s.userService.GetByID(userId)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting user",
			"data":    nil,
		})
		return
	}
	user.Address.ZipCode = addressInput.ZipCode
	user.Address.City = addressInput.City
	user.Address.Street = addressInput.Street
	user.Address.State = addressInput.State
	user.Address.HouseNumber = addressInput.HouseNumber
	user.Address.Country = addressInput.Country
	err = s.addressService.UpdateAddress(&user.Address)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error updating address",
			"data":    nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Address updated ",
		"data":    user,
	})
}
