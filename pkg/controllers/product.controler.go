package controllers

import (
	"Ecommerce/pkg/db"
	"Ecommerce/pkg/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// /api/v1/products/create
func CreateProduct(c *gin.Context) {
	session := db.GetDB().Session(&gorm.Session{})
	userType, ok := c.Get("userType")
	if !ok || userType != "ADMIN" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Only admin can get users",
			"data":    nil,
		})
		return
	}
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid product data",
			"data":    err.Error(),
		})
		return
	}
	result := session.Create(&product)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error creating product",
			"data":    nil,
		})
		return
	}
	session.Commit()
	c.JSON(201, gin.H{
		"status":  "success",
		"message": "Product created",
		"data":    product,
	})
}

// / /api/v1/products/
func GetAllProducts(c *gin.Context) {
	session := db.GetDB().Session(&gorm.Session{})
	var products []models.Product
	if err := session.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error getting products GetAllProducts",
		})
		return
	}
	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "No products found GetAllProducts",
			"data":    nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Users found GetAllProducts",
		"data":    products,
	})
}

// /api/v1/products/:id +
func GetProduct(c *gin.Context) {
	session := db.GetDB().Session(&gorm.Session{})
	id := c.Param("id")
	var product models.Product
	if err := session.Where("ID=?", id).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Product not found GetProduct",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Error getting product GetProduct",
			})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product found GetProduct",
		"data":    product,
	})
}

// /api/v1/products/:id
func UpdateProduct(c *gin.Context) {
	session := db.GetDB().Session(&gorm.Session{})
	userType, ok := c.Get("userType")
	if !ok || userType != "ADMIN" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Only admin can update products UpdateProduct",
			"data":    nil,
		})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid product id UpdateProduct",
		})
		return
	}
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid product data UpdateProduct",
			"data":    err.Error(),
		})
		return
	}
	result := session.Model(&models.Product{}).Where("id = ?", id).Updates(product)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error updating product UpdateProduct",
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product updated UpdateProduct",
		"data":    product,
	})
}

// /api/v1/products/:id
func DeleteProduct(c *gin.Context) {
	session := db.GetDB().Session(&gorm.Session{})
	userType, ok := c.Get("userType")
	if !ok || userType != "ADMIN" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Only admin can delete products DeleteProduct",
		})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid product id DeleteProduct",
		})
		return
	}
	var product models.Product
	res := session.Where("ID = ?", id).First(&product)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Product not found DeleteProduct",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error finding product DeleteProduct",
		})
		return
	}
	deleteResult := session.Delete(&product)
	if deleteResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error deleting product DeleteProduct",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Product deleted successfully DeleteProduct",
	})
}
