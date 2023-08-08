package controllers

import (
	"Ecommerce/pkg/db"
	"Ecommerce/pkg/models"
	"github.com/gin-gonic/gin"
)

///по любому не правильные

func RemoveProductFromCart(c *gin.Context) {
	type Request struct {
		Count int `json:"count" binding:"required"`
	}

	// parse req body
	var req Request
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid request body",
			"data":    err,
		})
		return
	}

	// get product id from params and user id from context
	productId := c.Param("id")
	if productId == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid product id",
			"data":    nil,
		})
		return
	}

	var product models.Product
	res := db.GetDB().Where("ID=?", productId).Find(&product)
	if res.Error != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Product with such an ID was not found",
			"data":    res.Error,
		})
		return
	}
	idLocal, ok := c.Get("user_id")
	if !ok {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "User ID not found in context",
		})
		return
	}
	var user models.User
	result := db.GetDB().Where("ID = ?", idLocal).Preload("UserCart").Find(&user)
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
		if string(item.ID) == productId {
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

	if err := db.GetDB().Save(&user).Error; err != nil {
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

func AddProductToCart(c *gin.Context) {
	idUser, ok := c.Get("id")
	if !ok {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "User ID not found in context",
		})
		return
	}
	idProduct := c.Param("id")
	if idProduct == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid product id",
		})
		return
	}
	// find product by id
	var product models.Product
	res := db.GetDB().Where("ID=?", idProduct).Find(&product)
	if res.Error != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Product not found",
			"data":    res.Error,
		})
		return
	}
	// check if product is available
	if product.AvailableQuantity == 0 {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Product is not available",
		})
		return
	}
	// fetch UserCart from user
	var user models.User
	result := db.GetDB().Where("ID = ?", idUser).Preload("UserCart").Find(&user)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "User with such an ID was not found",
			"data":    result.Error,
		})
		return
	}
	// parse req body
	var productToOrder models.ProductsToOrder
	if err := c.ShouldBindJSON(&productToOrder); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid request body",
			"data":    err.Error(),
		})
		return
	}
	// user quantity must be less than product quantity
	if productToOrder.BuyQuantity > product.AvailableQuantity {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Quantity must be less than product quantity",
			"data":    nil,
		})
		return
	}
	var existingCartItem models.ProductsToOrder
	existingIndex := -1
	for i, item := range user.UserCart {
		if string(item.ID) == idProduct {
			existingCartItem = item
			existingIndex = i
			break
		}
	}

	if existingIndex != -1 {
		user.UserCart[existingIndex].BuyQuantity += productToOrder.BuyQuantity
	} else {
		// Add the product to the user's cart
		productToOrder.Name = product.Name
		productToOrder.Price = product.Price
		productToOrder.ID = product.ID // Assuming you have a ProductID field
		user.UserCart = append(user.UserCart, productToOrder)
	}

	// Save the updated user cart
	if err := db.GetDB().Save(&user).Error; err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to update user cart",
			"data":    err,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Product added to cart successfully",
	})
}
