package controllers

import (
	"Ecommerce/pkg/db"
	"Ecommerce/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

// /api/v1/admin/getUser/:id
func GetUser(c *gin.Context) {
	session := db.GetDB().Session(&gorm.Session{})
	id := c.Param("id")
	ID, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid user id",
		})
		return
	}
	var user models.User
	res := session.Where("ID=?", ID).First(&user)
	if res.Error != nil {
		c.JSON(404, gin.H{
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

// /api/v1/admin/getUsers
func GetUsers(c *gin.Context) {
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
	var users []models.User
	res := session.Find(&users)
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

// /api/v1/admin/deleteUser/:id
func DeleteUser(c *gin.Context) {
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
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid user id",
		})
		return
	}
	var user models.User
	res := session.Where("ID=?", id).Delete(&user)
	if res.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error deleting users",
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

// /api/v1/admin/deleteUsers
func DeleteAllUsers(c *gin.Context) {
	session := db.GetDB().Session(&gorm.Session{})
	userType, ok := c.Get("userType")
	if !ok || userType != "ADMIN" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Only admin can get users DeleteAllUsers",
			"data":    nil,
		})
		return
	}
	var users []models.User
	result := session.Find(&users)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting users DeleteAllUsers",
			"data":    nil,
		})
		return
	}
	if len(users) == 0 {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "No users found for deletion DeleteAllUsers",
			"data":    nil,
		})
		return
	}
	deleteResult := session.Delete(&users)
	if deleteResult.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error deleting users DeleteAllUsers",
			"data":    nil,
		})
		return
	}
	session.Commit()
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Users deleted DeleteAllUsers",
		"data":    nil,
	})
}
