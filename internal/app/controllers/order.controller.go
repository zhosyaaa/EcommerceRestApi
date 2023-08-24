package controllers

import (
	_interface2 "Ecommerce/internal/app/service/interface"
	"Ecommerce/internal/pkg/db/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type OrderController struct {
	orderService   _interface2.OrderRepository
	productService _interface2.ProductRepository
	userService    _interface2.UserRepository
}

func NewOrderController(orderService _interface2.OrderRepository, productService _interface2.ProductRepository, userService _interface2.UserRepository) *OrderController {
	return &OrderController{orderService: orderService, productService: productService, userService: userService}
}

func (s *OrderController) OrderAll(c *gin.Context) {
	userID, _ := c.Get("id")
	user, err := s.userService.GetByID(userID.(string))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "User does not exist",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Internal server error ",
			})
		}
		return
	}

	userCart := user.UserCart
	if len(userCart) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User cart is empty",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
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
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "User does not exist",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
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
		c.JSON(http.StatusNotFound, gin.H{
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update user's cart",
			"data":    err,
		})
		return
	}

	err = s.orderService.DeleteProductToOrder(&productToOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to deleted productToOrder",
			"data":    err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Order removed from cart successfully ",
	})
}
