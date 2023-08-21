package controllers

import (
	"Ecommerce/pkg/models"
	interfaces "Ecommerce/pkg/repository/interface"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

type OrderController struct {
	orderService   interfaces.OrderRepository
	productService interfaces.ProductRepository
	userService    interfaces.UserRepository
}

func NewOrderController(orderService interfaces.OrderRepository, productService interfaces.ProductRepository, userService interfaces.UserRepository) *OrderController {
	return &OrderController{orderService: orderService, productService: productService, userService: userService}
}

func (s *OrderController) OrderAll(c *gin.Context) {
	userID, _ := c.Get("id")
	user, err := s.userService.GetByID(userID.(string))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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

func (s *OrderController) OrderOne(c *gin.Context) {
	userID, _ := c.Get("id")
	user, err := s.userService.GetByID(userID.(string))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
	err = s.userService.UpdateUser(id, user)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to update user's cart",
			"data":    err,
		})
		return
	}

	err = s.orderService.DeleteProductToOrder(&productToOrder)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to deleted productToOrder",
			"data":    err,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Order removed from cart successfully ",
	})
}
