package controllers

import (
	interfaces "Ecommerce/internal/app/service/interface"
	"Ecommerce/internal/pkg/db/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type AdminController struct {
	userService interfaces.UserRepository
}

func NewAdminController(userService interfaces.UserRepository) *AdminController {
	return &AdminController{userService: userService}
}

func (s *AdminController) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := s.userService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User found",
		"data":    user,
	})
}

func (s *AdminController) GetUsers(c *gin.Context) {
	userType, ok := c.Get("userType")
	if !ok || userType != "ADMIN" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Only admin can get users",
			"data":    nil,
		})
		return
	}
	var users []models.User
	users, err := s.userService.GetAllUsers(users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error get users",
			"data":    nil,
		})
		return
	}
	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "No users found",
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Users found",
		"data":    users,
	})
}

func (s *AdminController) DeleteUser(c *gin.Context) {
	userType, ok := c.Get("userType")
	if !ok || userType != "ADMIN" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Only admin can get users",
			"data":    nil,
		})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid user id",
		})
		return
	}
	user, err := s.userService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User not found",
			"data":    nil,
		})
		return
	}
	err = s.userService.DeleteUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error deleting user",
			"data":    nil,
		})
		return
	}
	if user.UserType == "ADMIN" {
		cookie := http.Cookie{
			Name:     "jwt",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			HttpOnly: true,
		}
		http.SetCookie(c.Writer, &cookie)
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Admin deleted",
			"data":    nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "User deleted",
			"data":    nil,
		})
	}
}

func (s *AdminController) DeleteAllUsers(c *gin.Context) {
	userType, ok := c.Get("userType")
	if !ok || userType != "ADMIN" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Only admin can get users",
			"data":    nil,
		})
		return
	}
	var users []models.User
	users, err := s.userService.GetAllUsers(users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error getting users",
			"data":    nil,
		})
		return
	}
	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "No users found for deletion",
			"data":    nil,
		})
		return
	}
	err = s.userService.DeleteAllUsers(users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error deleting users",
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Users deleted",
		"data":    nil,
	})
}
