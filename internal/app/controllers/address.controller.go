package controllers

import (
	_interface2 "Ecommerce/internal/app/service/interface"
	"Ecommerce/internal/pkg/db/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AddressController struct {
	addressService _interface2.AddressRepository
	userService    _interface2.UserRepository
}

func NewAddressController(addressService _interface2.AddressRepository, userService _interface2.UserRepository) *AddressController {
	return &AddressController{addressService: addressService, userService: userService}
}

func (s *AddressController) UpdateAddress(c *gin.Context) {
	userId := c.Param("id")
	var addressInput models.AddressInputCred
	if err := c.ShouldBindJSON(&addressInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid address data",
			"data":    err.Error(),
		})
		return
	}
	user, err := s.userService.GetByID(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error updating address",
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Address updated ",
		"data":    user,
	})
}
