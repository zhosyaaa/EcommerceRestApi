package controllers

import (
	"Ecommerce/pkg/db"
	"Ecommerce/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
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
				"message": "User does not exist OrderAll",
			})
		} else {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "Internal server error OrderAll",
			})
		}
		return
	}

	userCart := user.UserCart
	if len(userCart) == 0 {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "User cart is empty OrderAll",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User cart retrieved successfully OrderAll32 OrderAll",
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
				"message": "User does not exist OrderOne",
			})
		} else {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "Internal server error OrderOne",
			})
		}
		return
	}

	id := c.Param("id")
	var productToOrder models.ProductsToOrder
	var cartIndex int
	for i, item := range user.UserCart {
		if strconv.Itoa(int(item.ID)) == id {
			productToOrder = item
			cartIndex = i
			break
		}
	}
	if productToOrder.ID == 0 {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Order does not exist OrderOne",
		})
		return
	}
	user.UserCart = append(user.UserCart[:cartIndex], user.UserCart[cartIndex+1:]...)
	updateResult := db.GetDB().Save(&user)
	if updateResult.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to update user's cart OrderOne",
			"data":    updateResult.Error,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Order removed from cart successfully OrderOne",
	})
}
