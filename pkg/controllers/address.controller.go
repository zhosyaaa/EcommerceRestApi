package controllers

import (
	"Ecommerce/pkg/db"
	"Ecommerce/pkg/models"
	"github.com/gin-gonic/gin"
)

func UpdateAddress(c *gin.Context) {
	userId := c.Param("id")
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid address data",
			"data":    err.Error(),
		})
		return
	}
	result := db.GetDB().Model(&models.User{}).Where("id = ?", userId).Updates(map[string]interface{}{"address": user.Address})
	if result.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error updating address",
			"data":    nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Address updated",
		"data":    user,
	})
}
