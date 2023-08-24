package controllers

import (
	_interface2 "Ecommerce/internal/app/service/interface"
	"Ecommerce/internal/pkg/db/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CartController struct {
	orderService   _interface2.OrderRepository
	productService _interface2.ProductRepository
	userService    _interface2.UserRepository
}

func NewCartController(orderService _interface2.OrderRepository, productService _interface2.ProductRepository, userService _interface2.UserRepository) *CartController {
	return &CartController{orderService: orderService, productService: productService, userService: userService}
}

func (s *CartController) RemoveProductFromCart(c *gin.Context) {
	type Request struct {
		Count int `json:"count" binding:"required"`
	}
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body",
			"data":    err,
		})
		return
	}
	productId := c.Param("id")
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid product id ",
			"data":    nil,
		})
		return
	}
	userID, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User ID not found ",
		})
		return
	}
	product, err := s.productService.GetByID(productId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Product with such an ID was not found",
			"data":    err,
		})
		return
	}
	var user *models.User
	user, err = s.userService.GetByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User with such an ID was not found",
			"data":    err,
		})
		return
	}
	var found bool
	var PO models.ProductsToOrder
	for i, item := range user.UserCart {
		if strconv.Itoa(int(item.ID)) == productId {
			if item.BuyQuantity >= req.Count {
				item.BuyQuantity -= req.Count
			} else {
				user.UserCart = append(user.UserCart[:i], user.UserCart[i+1:]...)
			}
			found = true
			product.AvailableQuantity += req.Count
			PO = item
			break
		}
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Product not found in cart",
		})
		return
	}
	err = s.productService.UpdateProduct(productId, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error updating product",
			"data":    nil,
		})
		return
	}
	err = s.orderService.UpdateProductToOrder(productId, &PO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error updating ProductsToOrder",
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product removed from cart successfully",
	})
}

func (s *CartController) AddProductToCart(c *gin.Context) {
	idUser, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User ID not found",
		})
		return
	}
	idProduct := c.Param("id")
	if idProduct == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid product ID",
		})
		return
	}

	product, err := s.productService.GetByID(idProduct)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Product not found ",
			"data":    err.Error(),
		})
		return
	}
	if product.AvailableQuantity == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Product is not available",
		})
		return
	}

	user, err := s.userService.GetByID(idUser.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User with such an ID was not found",
			"data":    err.Error(),
		})
		return
	}

	var productToOrder *models.ProductsToOrder
	if err := c.ShouldBindJSON(&productToOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body",
			"data":    err.Error(),
		})
		return
	}
	if productToOrder.BuyQuantity > product.AvailableQuantity {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Quantity must be less than product quantity",
		})
		return
	}

	var orders models.Order
	orders.TotalPrice = float64(productToOrder.Price) * float64(productToOrder.BuyQuantity)
	orders.AddressID = user.ID
	orders.OrderCart = append(orders.OrderCart, *productToOrder)

	err = s.orderService.CreateOrder(&orders)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to insert a order",
			"data":    err,
		})
		return
	}
	user.Orders = append(user.Orders, orders)
	if len(user.UserCart) == 0 {
		user.UserCart = append(user.UserCart, *productToOrder)
	} else {
		existingIndex := -1
		for i, item := range user.UserCart {
			if item.ID == product.ID {
				existingIndex = i
				break
			}
		}
		updateCart(user, *product, *productToOrder, existingIndex)
	}
	subtractQuantity := product.AvailableQuantity - productToOrder.BuyQuantity

	product.AvailableQuantity = subtractQuantity
	err = s.productService.UpdateProduct(strconv.Itoa(int(product.ID)), product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error updating product",
			"data":    nil,
		})
		return
	}

	err = s.userService.UpdateUser(strconv.Itoa(int(user.ID)), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error updating user",
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product added to cart successfully",
		"data":    productToOrder,
	})
}

func updateCart(user *models.User, product models.Product, productToOrder models.ProductsToOrder, existingIndex int) {
	if existingIndex != -1 {
		user.UserCart[existingIndex].BuyQuantity += productToOrder.BuyQuantity
	} else {
		productToOrder.Name = product.Name
		productToOrder.Price = product.Price
		productToOrder.ID = product.ID
		user.UserCart = append(user.UserCart, productToOrder)
	}
}
