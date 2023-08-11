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
	userID, _ := c.Get("id")
	var user models.User
	result := db.GetDB().Where("ID=?", userID).Preload("UserCart").First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"status":  "error",
				"message": "User does not exist",
			})
		} else {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "Internal server error ",
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

// /api/v1/order/:id +
func OrderOne(c *gin.Context) {
	session := db.GetDB().Session(&gorm.Session{})
	userID, _ := c.Get("id")
	var user models.User
	result := session.Where("ID=?", userID).Preload("UserCart").First(&user)
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
		if strconv.Itoa(int(item.ID)) == id {
			productToOrder = item
			cartIndex = i
			break
		}
	}
	if productToOrder.ID == 0 {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Order does not exist",
		})
		return
	}
	var usercart []models.ProductsToOrder
	usercart = append(user.UserCart[:cartIndex], user.UserCart[cartIndex+1:]...)
	user.UserCart = usercart
	updateResult := session.Save(&user)
	if updateResult.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to update user's cart",
			"data":    updateResult.Error,
		})
		return
	}

	updateResult = session.Delete(&productToOrder)
	if updateResult.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to deleted productToOrder",
			"data":    updateResult.Error,
		})
		return
	}
	session.Commit()
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Order removed from cart successfully ",
	})
}
