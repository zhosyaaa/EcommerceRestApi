package controllers

import (
	"Ecommerce/pkg/db"
	"Ecommerce/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// /api/v1/address/update/:id
func UpdateAddress(c *gin.Context) {
	session := db.GetDB().Session(&gorm.Session{})
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
	result := session.Model(&models.User{}).Where("id = ?", userId).Updates(map[string]interface{}{"address": user.Address})
	if result.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error updating address",
			"data":    nil,
		})
		return
	}
	session.Commit()
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Address updated",
		"data":    user,
	})
}
