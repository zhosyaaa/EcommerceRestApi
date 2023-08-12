package controllers

import (
	"Ecommerce/pkg/db"
	"Ecommerce/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

// /api/v1/admin/getUser/:id
func GetUser(c *gin.Context) {
	logger := log.With().Str("request_id", c.GetString("x-request-id")).Logger()
	logger.Debug().Msg("Received request to Get User")

	session := db.GetDB().Session(&gorm.Session{})
	id := c.Param("id")
	ID, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		logger.Error().Err(err).Msg("Invalid user id")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid user id",
		})
		return
	}
	var user models.User
	res := session.Where("ID=?", ID).Preload("Address").Preload("UserCart").Preload("Orders").First(&user)
	if res.Error != nil {
		logger.Error().Err(res.Error).Msg("User not found")
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}
	logger.Info().Int64("user_id", int64(user.ID)).Msg("User found")
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User found",
		"data":    user,
	})
}

// /api/v1/admin/getUsers
func GetUsers(c *gin.Context) {
	logger := log.With().Str("request_id", c.GetString("x-request-id")).Logger()
	logger.Debug().Msg("Received request to Remove Get Users")

	session := db.GetDB().Session(&gorm.Session{})
	userType, ok := c.Get("userType")
	if !ok || userType != "ADMIN" {
		logger.Warn().Msg("Only admin can get users")
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Only admin can get users",
			"data":    nil,
		})
		return
	}
	var users []models.User
	res := session.Preload("Address").Preload("UserCart").Preload("Orders").Find(&users)
	if res.Error != nil {
		logger.Error().Err(res.Error).Msg("Error getting users")
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error get users",
			"data":    nil,
		})
		return
	}
	if len(users) == 0 {
		logger.Warn().Msg("No users found")
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "No users found",
			"data":    nil,
		})
		return
	}
	logger.Info().Msg("Users found")
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Users found",
		"data":    users,
	})
}

// /api/v1/admin/deleteUser/:id
func DeleteUser(c *gin.Context) {
	logger := log.With().Str("request_id", c.GetString("x-request-id")).Logger()
	logger.Debug().Msg("Received request to Delete User")

	session := db.GetDB().Session(&gorm.Session{})
	userType, ok := c.Get("userType")
	if !ok || userType != "ADMIN" {
		logger.Warn().Msg("Only admin can delete users")
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Only admin can get users",
			"data":    nil,
		})
		return
	}
	id := c.Param("id")
	if id == "" {
		logger.Warn().Msg("Invalid user id")
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid user id",
		})
		return
	}
	var user models.User
	res := session.Where("ID=?", id).Preload("Address").Preload("UserCart").Preload("Orders").First(&user)
	if res.Error != nil {
		logger.Error().Err(res.Error).Msg("Error getting user")
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting user",
			"data":    nil,
		})
		return
	}

	session.Delete(&user.Address)
	session.Delete(&user.UserCart)
	session.Delete(&user.Orders)

	res = session.Delete(&user)
	if res.Error != nil {
		logger.Error().Err(res.Error).Msg("Error deleting user")
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error deleting user",
			"data":    nil,
		})
		return
	}

	session.Commit()
	if user.UserType == "ADMIN" {
		cookie := http.Cookie{
			Name:     "jwt",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			HttpOnly: true,
		}
		http.SetCookie(c.Writer, &cookie)
		logger.Info().Msg("Admin deleted")
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Admin deleted",
			"data":    nil,
		})
	} else {
		logger.Info().Msg("User deleted")
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "User deleted",
			"data":    nil,
		})
	}
}

// /api/v1/admin/deleteUsers
func DeleteAllUsers(c *gin.Context) {
	logger := log.With().Str("request_id", c.GetString("x-request-id")).Logger()
	logger.Debug().Msg("Received request to Delete All Users")

	session := db.GetDB().Session(&gorm.Session{})
	userType, ok := c.Get("userType")
	if !ok || userType != "ADMIN" {
		logger.Warn().Msg("Only admin can delete users")
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Only admin can get users",
			"data":    nil,
		})
		return
	}
	var users []models.User
	result := session.Find(&users)
	if result.Error != nil {
		logger.Error().Err(result.Error).Msg("Error getting users")
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting users",
			"data":    nil,
		})
		return
	}
	if len(users) == 0 {
		logger.Warn().Msg("No users found for deletion")
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "No users found for deletion",
			"data":    nil,
		})
		return
	}

	for _, user := range users {
		session.Delete(&user.Address)
		session.Delete(&user.UserCart)
		session.Delete(&user.Orders)
	}

	deleteResult := session.Delete(&users)
	if deleteResult.Error != nil {
		logger.Error().Err(deleteResult.Error).Msg("Error deleting users")
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error deleting users",
			"data":    nil,
		})
		return
	}
	session.Commit()
	logger.Info().Msg("Users deleted")
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Users deleted",
		"data":    nil,
	})
}
