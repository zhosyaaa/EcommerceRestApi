package controllers

import (
	"Ecommerce/pkg/db"
	"Ecommerce/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"strconv"
)

// /api/v1/order/
func OrderAll(c *gin.Context) {
	userID, _ := c.Get("id")
	logger := log.With().Str("request_id", c.GetString("x-request-id")).Logger()
	logger.Debug().Msg("Received request to retrieve user cart")

	var user models.User
	result := db.GetDB().Where("ID=?", userID).Preload("UserCart").First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.Warn().Msg("User not found")
			c.JSON(404, gin.H{
				"status":  "error",
				"message": "User does not exist",
			})
		} else {
			logger.Error().Err(result.Error).Msg("Internal server error")
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "Internal server error ",
			})
		}
		return
	}

	userCart := user.UserCart
	if len(userCart) == 0 {
		logger.Warn().Msg("User cart is empty")
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "User cart is empty",
		})
		return
	}

	logger.Info().Int("num_items", len(userCart)).Msg("User cart retrieved successfully")
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
	logger := log.With().Str("request_id", c.GetString("x-request-id")).Logger()
	logger.Debug().Msg("Received request to remove order from user cart")

	var user models.User
	result := session.Where("ID=?", userID).Preload("UserCart").First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.Warn().Msg("User not found")
			c.JSON(404, gin.H{
				"status":  "error",
				"message": "User does not exist",
			})
		} else {
			logger.Error().Err(result.Error).Msg("Internal server error")
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
		logger.Warn().Msg("Order not found")
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
		logger.Error().Err(updateResult.Error).Msg("Failed to update user's cart")
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to update user's cart",
			"data":    updateResult.Error,
		})
		return
	}

	updateResult = session.Delete(&productToOrder)
	if updateResult.Error != nil {
		logger.Error().Err(updateResult.Error).Msg("Failed to delete productToOrder")
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to deleted productToOrder",
			"data":    updateResult.Error,
		})
		return
	}
	session.Commit()
	logger.Info().Int("product_id", int(productToOrder.ID)).Msg("Order removed from cart successfully")
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Order removed from cart successfully ",
	})
}
