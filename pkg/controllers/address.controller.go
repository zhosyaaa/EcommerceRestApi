package controllers

import (
	"Ecommerce/pkg/db"
	"Ecommerce/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// /api/v1/address/update/:id
func UpdateAddress(c *gin.Context) {
	logger := log.With().Str("request_id", c.GetString("x-request-id")).Logger()
	logger.Debug().Msg("Received request to Update Address")

	session := db.GetDB().Session(&gorm.Session{})
	userId := c.Param("id")
	type UserInputCred struct {
		ZipCode     string `json:"zipCode"`
		City        string `json:"city"`
		State       string `json:"state"`
		Country     string `json:"country"`
		Street      string `json:"street"`
		HouseNumber string `json:"houseNumber"`
	}
	var addressInput UserInputCred
	if err := c.ShouldBindJSON(&addressInput); err != nil {
		logger.Error().Err(err).Msg("Invalid address data")
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid address data",
			"data":    err.Error(),
		})
		return
	}
	var user models.User
	result := session.Where("id = ?", userId).Preload("Address").First(&user) // Added '&' before user
	if result.Error != nil {
		logger.Error().Err(result.Error).Msg("Error getting user")
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
	res := session.Save(&user.Address)
	if res.Error != nil {
		logger.Error().Err(res.Error).Msg("Error updating address")
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error updating address",
			"data":    nil,
		})
		return
	}
	session.Commit()
	logger.Info().Msg("Address updated")
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Address updated ",
		"data":    user,
	})
}
