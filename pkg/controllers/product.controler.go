package controllers

import (
	"Ecommerce/pkg/db"
	"Ecommerce/pkg/models"
	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	userType := c.Param("userType")
	if userType == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid user data for model binding",
			"data":    nil,
		})
		return
	}
	if userType != "ADMIN" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Only admin can create products",
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

	result := db.GetDB().Create(&product)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error creating product",
			"data":    nil,
		})
		return
	}

	c.JSON(201, gin.H{
		"status":  "success",
		"message": "Product created",
		"data":    product,
	})
}
func GetAllProducts(c *gin.Context) {
	var products []models.Product
	res := db.GetDB().Find(&products)
	if res.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting products",
			"data":    nil,
		})
		return
	}
	if len(products) == 0 {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "No products found",
			"data":    nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Users found",
		"data":    products,
	})
}
func GetProduct(c *gin.Context) {
	id := c.Param("id")
	var product *models.Product
	res := db.GetDB().Where("ID=?", id).Find(&product)
	if res.Error != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "No product found",
			"data":    nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Users found",
		"data":    product,
	})
}

func UpdateProduct(c *gin.Context) {
	userType := c.Param("userType")
	if userType == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid user data for model binding",
			"data":    nil,
		})
		return
	}
	if userType != "ADMIN" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Only admin can create products",
			"data":    nil,
		})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid product id",
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
	result := db.GetDB().Model(&models.Product{}).Where("id = ?", id).Updates(product)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error updating product",
			"data":    nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Product updated",
		"data":    product,
	})
}
func DeleteProduct(c *gin.Context) {
	userType := c.Param("userType")
	if userType == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid user data for model binding",
			"data":    nil,
		})
		return
	}
	if userType != "ADMIN" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Only admin can create products",
			"data":    nil,
		})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid product id",
		})
		return
	}
	var product models.Product
	res := db.GetDB().Where("ID=?", id).Delete(&product)
	if res.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error deleting users",
			"data":    nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Product deleted successfully",
		"data":    nil,
	})
}
