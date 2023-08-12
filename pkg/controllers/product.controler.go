package controllers

import (
	"Ecommerce/pkg/db"
	"Ecommerce/pkg/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
)

// /api/v1/products/create
func CreateProduct(c *gin.Context) {
	logger := log.With().Str("request_id", c.GetString("x-request-id")).Logger()
	logger.Debug().Msg("Received request to create product")

	session := db.GetDB().Session(&gorm.Session{})
	userType, ok := c.Get("userType")
	if !ok || userType != "ADMIN" {
		logger.Warn().Msg("Unauthorized request: Only admin can create products")
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Only admin can get users",
			"data":    nil,
		})
		return
	}
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		logger.Error().Err(err).Msg("Invalid product data")
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid product data",
			"data":    err.Error(),
		})
		return
	}
	result := session.Create(&product)
	if result.Error != nil {
		logger.Error().Err(result.Error).Msg("Error creating product")
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error creating product",
			"data":    nil,
		})
		return
	}
	session.Commit()
	logger.Info().Int("product_id", int(product.ID)).Msg("Product created successfully")
	c.JSON(201, gin.H{
		"status":  "success",
		"message": "Product created",
		"data":    product,
	})
}

// / /api/v1/products/
func GetAllProducts(c *gin.Context) {
	session := db.GetDB().Session(&gorm.Session{})
	logger := log.With().Str("request_id", c.GetString("x-request-id")).Logger()
	logger.Debug().Msg("Received request to get all products")

	var products []models.Product
	if err := session.Find(&products).Error; err != nil {
		logger.Error().Err(err).Msg("Error getting products")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error getting products",
		})
		return
	}
	if len(products) == 0 {
		logger.Warn().Msg("No products found")
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "No products found",
			"data":    nil,
		})
		return
	}

	logger.Info().Int("num_products", len(products)).Msg("Products found")
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Users found",
		"data":    products,
	})
}

// /api/v1/products/:id +
func GetProduct(c *gin.Context) {
	session := db.GetDB().Session(&gorm.Session{})
	id := c.Param("id")

	logger := log.With().Str("request_id", c.GetString("x-request-id")).Logger()
	logger.Debug().Str("product_id", id).Msg("Received request to get product")

	var product models.Product
	if err := session.Where("ID=?", id).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn().Msg("Product not found")
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Product not found",
			})
		} else {
			logger.Error().Err(err).Msg("Error getting product")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Error getting product",
			})
		}
		return
	}
	logger.Info().Str("product_name", product.Name).Msg("Product found")
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product found",
		"data":    product,
	})
}

// /api/v1/products/:id
func UpdateProduct(c *gin.Context) {
	logger := log.With().Str("request_id", c.GetString("x-request-id")).Logger()
	logger.Debug().Msg("Received request to update product")

	session := db.GetDB().Session(&gorm.Session{})
	userType, ok := c.Get("userType")
	if !ok || userType != "ADMIN" {
		logger.Warn().Msg("Unauthorized request: Only admin can update products")
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Only admin can update products",
			"data":    nil,
		})
		return
	}
	id := c.Param("id")
	if id == "" {
		logger.Warn().Msg("Invalid product id")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid product id",
		})
		return
	}
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		logger.Error().Err(err).Msg("Invalid product data")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid product data",
			"data":    err.Error(),
		})
		return
	}
	result := session.Model(&models.Product{}).Where("id = ?", id).Updates(product)
	if result.Error != nil {
		logger.Error().Err(result.Error).Msg("Error updating product")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error updating product",
			"data":    nil,
		})
		return
	}
	logger.Info().Int("product_id", int(product.ID)).Msg("Product updated successfully")
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product updated",
		"data":    product,
	})
}

// /api/v1/products/:id
func DeleteProduct(c *gin.Context) {
	logger := log.With().Str("request_id", c.GetString("x-request-id")).Logger()
	logger.Debug().Msg("Received request to delete product")

	session := db.GetDB().Session(&gorm.Session{})
	userType, ok := c.Get("userType")
	if !ok || userType != "ADMIN" {
		logger.Warn().Msg("Unauthorized request: Only admin can delete products")
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Only admin can delete products",
		})
		return
	}
	id := c.Param("id")
	if id == "" {
		logger.Warn().Msg("Invalid product id")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid product id",
		})
		return
	}

	var product models.Product
	res := session.Where("ID = ?", id).First(&product)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			logger.Warn().Msg("Product not found")
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Product not found",
			})
			return
		}
		logger.Error().Err(res.Error).Msg("Error finding product")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error finding product",
		})
		return
	}
	deleteResult := session.Delete(&product)
	if deleteResult.Error != nil {
		logger.Error().Err(deleteResult.Error).Msg("Error deleting product")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error deleting product",
		})
		return
	}
	logger.Info().Int("product_id", int(product.ID)).Msg("Product deleted successfully")
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Product deleted successfully",
	})
}
