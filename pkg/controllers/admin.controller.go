package controllers

import (
	"Ecommerce/pkg/db"
	"Ecommerce/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetUser(c *gin.Context) {
	database := db.GetDB()
	userType := c.Param("userType")
	if userType == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid user data for model binding",
			"data":    nil,
		})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid user id",
		})
		return
	}
	// find user by id
	var user models.User
	res := database.First(&user, id)
	if res.Error != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User found",
		"data":    user,
	})
}

func GetUsers(c *gin.Context) {
	database := db.GetDB()
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
			"message": "Only admin can get users",
			"data":    nil,
		})
		return
	}

	var users []models.User
	res := database.Find(&users)
	if res.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error get users",
			"data":    nil,
		})
		return
	}
	if len(users) == 0 {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "No users found",
			"data":    nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Users found",
		"data":    users,
	})
}

func DeleteUser(c *gin.Context) {
	database := db.GetDB()
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
			"message": "Only admin can get users",
			"data":    nil,
		})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid user id",
		})
		return
	}
	var user models.User
	res := database.Where("ID=?", id).Delete(&user)
	if res.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error deleting users",
			"data":    nil,
		})
		return
	}
	if user.UserType == "ADMIN" {
		// delete cookie
		cookie := http.Cookie{
			Name:     "jwt",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			HttpOnly: true,
		}
		http.SetCookie(c.Writer, &cookie)
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Admin deleted",
			"data":    nil,
		})
	} else {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "User deleted",
			"data":    nil,
		})
	}
}

func DeleteAllUsers(c *gin.Context) {
	database := db.GetDB()
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
			"message": "Only admin can get users",
			"data":    nil,
		})
		return
	}
	var users []models.User

	result := database.Find(&users)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting users",
			"data":    nil,
		})
		return
	}
	if len(users) == 0 {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "No users found for deletion",
			"data":    nil,
		})
		return
	}
	deleteResult := database.Delete(&users)
	if deleteResult.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error deleting users",
			"data":    nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Users deleted",
		"data":    nil,
	})
}
