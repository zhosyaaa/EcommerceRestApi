package controllers

import (
	"Ecommerce/pkg/db"
	"Ecommerce/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// /api/v1/cart/remove/:id
func RemoveProductFromCart(c *gin.Context) {
	type Request struct {
		Count int `json:"count" binding:"required"`
	}
	var req Request
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid request body",
			"data":    err,
		})
		return
	}
	productId := c.Param("id")
	if productId == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid product id",
			"data":    nil,
		})
		return
	}
	userID, ok := c.Get("user_id") // _id
	if !ok {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "User ID not found in context",
		})
		return
	}
	session := db.GetDB().Session(&gorm.Session{})
	var product models.Product
	res := session.Where("ID=?", productId).Find(&product) // or first
	if res.Error != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Product with such an ID was not found",
			"data":    res.Error,
		})
		return
	}
	var user models.User
	result := session.Where("ID = ?", userID).Preload("UserCart").Find(&user) // or first
	if result.Error != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "User with such an ID was not found",
			"data":    result.Error,
		})
		return
	}

	var found bool
	for i, item := range user.UserCart {
		if string(item.ProductId) == productId {
			if user.UserCart[i].BuyQuantity >= req.Count {
				user.UserCart[i].BuyQuantity -= req.Count
			} else {
				user.UserCart = append(user.UserCart[:i], user.UserCart[i+1:]...)
			}
			found = true
			break
		}
	}
	if !found {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Product not found in cart",
		})
		return
	}

	if err := session.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to update user cart",
			"data":    err,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Product removed from cart successfully",
	})
}

// /api/v1/cart/add/:id
func AddProductToCart(c *gin.Context) {
	session := db.GetDB().Session(&gorm.Session{})
	idUser, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User ID not found in context",
		})
		return
	}
	userID, ok := idUser.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid user ID type",
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

	var product models.Product
	res := session.Where("ID=?", idProduct).Find(&product)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Product not found",
			"data":    res.Error,
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
	var user models.User
	result := session.Where("ID = ?", userID).Preload("UserCart").Find(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User with such an ID was not found",
			"data":    result.Error,
		})
		return
	}
	var productToOrder models.ProductsToOrder
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
	existingIndex := -1
	var existingCartItem models.ProductsToOrder
	for i, item := range user.UserCart {
		if item.ProductId == product.ID {
			existingCartItem = item
			existingIndex = i
			break
		}
	}
	updateCart(&user, product, productToOrder, existingCartItem, existingIndex)
	if err := session.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update user cart",
			"data":    err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product added to cart successfully",
	})
}

func updateCart(user *models.User, product models.Product, productToOrder models.ProductsToOrder, existingCartItem models.ProductsToOrder, existingIndex int) {
	if existingIndex != -1 {
		user.UserCart[existingIndex].BuyQuantity += productToOrder.BuyQuantity
	} else {
		productToOrder.Name = product.Name
		productToOrder.Price = product.Price
		productToOrder.ProductId = product.ID
		user.UserCart = append(user.UserCart, productToOrder)
	}
}
