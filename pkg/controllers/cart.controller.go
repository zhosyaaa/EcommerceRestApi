package controllers

import (
	"Ecommerce/pkg/db"
	"Ecommerce/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// /api/v1/cart/remove/:id
func RemoveProductFromCart(c *gin.Context) {
	logger := log.With().Str("request_id", c.GetString("x-request-id")).Logger()
	logger.Debug().Msg("Received request to Remove Product From Cart")

	type Request struct {
		Count int `json:"count" binding:"required"`
	}
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error().Err(err).Msg("Invalid request body")
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid request body",
			"data":    err,
		})
		return
	}
	productId := c.Param("id")
	if productId == "" {
		logger.Warn().Msg("Invalid product id")
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid product id ",
			"data":    nil,
		})
		return
	}
	userID, ok := c.Get("id")
	if !ok {
		logger.Warn().Msg("User ID not found")
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "User ID not found ",
		})
		return
	}
	session := db.GetDB().Session(&gorm.Session{})
	var product models.Product
	res := session.Where("ID=?", productId).Find(&product)
	if res.Error != nil {
		logger.Error().Err(res.Error).Msg("Product with such an ID was not found")
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Product with such an ID was not found",
			"data":    res.Error,
		})
		return
	}
	var user models.User
	result := session.Where("ID = ?", userID).Preload("UserCart").Find(&user)
	if result.Error != nil {
		logger.Error().Err(result.Error).Msg("User with such an ID was not found")
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "User with such an ID was not found",
			"data":    result.Error,
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
		logger.Warn().Msg("Product not found in cart")
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Product not found in cart",
		})
		return
	}
	ress := session.Model(&models.Product{}).Where("id = ?", productId).Updates(product)
	if ress.Error != nil {
		session.Rollback()
		logger.Error().Err(ress.Error).Msg("Error updating product")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error updating product",
			"data":    nil,
		})
		return
	}
	ress = session.Model(&models.ProductsToOrder{}).Where("id = ?", productId).Updates(PO)
	if ress.Error != nil {
		session.Rollback()
		logger.Error().Err(ress.Error).Msg("Error updating ProductsToOrder")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error updating ProductsToOrder",
			"data":    nil,
		})
		return
	}

	session.Commit()
	logger.Info().Int("product_id", int(product.ID)).Msg("Product removed from cart successfully")
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Product removed from cart successfully",
	})
}

// /api/v1/cart/add/:id +
func AddProductToCart(c *gin.Context) {
	logger := log.With().Str("request_id", c.GetString("x-request-id")).Logger()
	logger.Debug().Msg("Received request to Remove Product From Cart")

	session := db.GetDB().Session(&gorm.Session{})
	idUser, ok := c.Get("id")
	if !ok {
		logger.Warn().Msg("User ID not found")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User ID not found",
		})
		return
	}
	idProduct := c.Param("id")
	if idProduct == "" {
		logger.Warn().Msg("Invalid product ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid product ID",
		})
		return
	}

	var product models.Product
	if err := session.First(&product, idProduct).Error; err != nil {
		logger.Error().Err(err).Msg("Product not found")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Product not found ",
			"data":    err.Error(),
		})
		return
	}
	if product.AvailableQuantity == 0 {
		logger.Warn().Msg("Product is not available")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Product is not available",
		})
		return
	}

	var user models.User
	if err := session.Preload("UserCart").First(&user, idUser).Error; err != nil {
		logger.Error().Err(err).Msg("User with such an ID was not found")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User with such an ID was not found",
			"data":    err.Error(),
		})
		return
	}

	var productToOrder models.ProductsToOrder
	if err := c.ShouldBindJSON(&productToOrder); err != nil {
		logger.Error().Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body",
			"data":    err.Error(),
		})
		return
	}
	if productToOrder.BuyQuantity > product.AvailableQuantity {
		logger.Warn().Msg("Quantity must be less than product quantity")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Quantity must be less than product quantity",
		})
		return
	}
	var orders models.Order
	orders.TotalPrice = float64(productToOrder.Price) * float64(productToOrder.BuyQuantity)
	orders.AddressID = user.ID
	orders.OrderCart = append(orders.OrderCart, productToOrder)
	result := session.Create(&orders)
	if result.Error != nil {
		logger.Error().Err(result.Error).Msg("Failed to insert an order")
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to insert a order",
			"data":    result.Error.Error(),
		})
		return
	}
	user.Orders = append(user.Orders, orders)
	if len(user.UserCart) == 0 {
		user.UserCart = append(user.UserCart, productToOrder)
	} else {
		existingIndex := -1
		for i, item := range user.UserCart {
			if item.ID == product.ID {
				existingIndex = i
				break
			}
		}
		updateCart(&user, product, productToOrder, existingIndex)
	}
	subtractQuantity := product.AvailableQuantity - productToOrder.BuyQuantity

	product.AvailableQuantity = subtractQuantity
	result = session.Model(&models.Product{}).Where("id = ?", product.ID).Updates(product)
	if result.Error != nil {
		logger.Error().Err(result.Error).Msg("Error updating product")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error updating product",
			"data":    nil,
		})
		return
	}

	ress := session.Model(&models.User{}).Where("id = ?", user.ID).Updates(user)
	if ress.Error != nil {
		logger.Error().Err(ress.Error).Msg("Error updating user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error updating user",
			"data":    nil,
		})
		return
	}
	logger.Info().Int("product_id", int(product.ID)).Msg("Product added to cart successfully")
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
