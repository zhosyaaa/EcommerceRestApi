package controllers

import (
	"Ecommerce/pkg/db"
	"Ecommerce/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// /api/v1/order/
func OrderAll(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user models.User
	result := db.GetDB().Preload("UserCart").Where("ID=?", userID).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"status":  "error",
				"message": "User does not exist",
			})
		} else {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "Internal server error",
			})
		}
		return
	}

	userCart := user.UserCart
	if len(userCart) == 0 {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "User cart is empty",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User cart retrieved successfully",
		"data":    userCart,
	})
}

// /api/v1/order/:id
func OrderOne(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user models.User
	result := db.GetDB().Preload("UserCart").Where("ID=?", userID).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"status":  "error",
				"message": "User does not exist",
			})
		} else {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "Internal server error",
			})
		}
		return
	}

	id := c.Param("id")
	var productToOrder models.ProductsToOrder
	var cartIndex int
	for i, item := range user.UserCart {
		if string(item.ProductId) == id {
			productToOrder = item
			cartIndex = i
			break
		}
	}
	if productToOrder.ProductId == 0 {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Order does not exist",
		})
		return
	}
	user.UserCart = append(user.UserCart[:cartIndex], user.UserCart[cartIndex+1:]...)
	updateResult := db.GetDB().Save(&user)
	if updateResult.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to update user's cart",
			"data":    updateResult.Error,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Order removed from cart successfully",
	})
}
