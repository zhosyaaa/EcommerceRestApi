package controllers

import (
	interfaces "Ecommerce/internal/app/service/interface"
	"Ecommerce/internal/pkg/db/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type ProductController struct {
	productService interfaces.ProductRepository
}

func NewProductController(productService interfaces.ProductRepository) *ProductController {
	return &ProductController{productService: productService}
}

func (s *ProductController) CreateProduct(c *gin.Context) {
	userType, ok := c.Get("userType")
	if !ok || userType != "ADMIN" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Only admin can get users",
			"data":    nil,
		})
		return
	}
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid product data",
			"data":    err.Error(),
		})
		return
	}
	err := s.productService.CreateProduct(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error creating product",
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Product created",
		"data":    product,
	})
}

func (s *ProductController) GetAllProducts(c *gin.Context) {
	var products []models.Product
	products, err := s.productService.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error getting products",
		})
		return
	}
	fmt.Println(products, products[0], len(products))
	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "No products found",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Products found",
		"data":    products,
	})
}

func (s *ProductController) GetProduct(c *gin.Context) {
	id := c.Param("id")
	ID, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid user id",
		})
		return
	}

	product, err := s.productService.GetByID(strconv.FormatInt(ID, 10))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Product not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Error getting product",
			})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product found",
		"data":    product,
	})
}

func (s *ProductController) UpdateProduct(c *gin.Context) {
	userType, ok := c.Get("userType")
	if !ok || userType != "ADMIN" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Only admin can update products",
			"data":    nil,
		})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid product id",
		})
		return
	}
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid product data",
			"data":    err.Error(),
		})
		return
	}
	err := s.productService.UpdateProduct(id, &product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error updating product",
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product updated",
		"data":    product,
	})
}

func (s *ProductController) DeleteProduct(c *gin.Context) {
	userType, ok := c.Get("userType")
	if !ok || userType != "ADMIN" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Only admin can delete products",
		})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid product id",
		})
		return
	}

	err := s.productService.DeleteProduct(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error deleting product",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product deleted successfully",
	})
}
