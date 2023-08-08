package controllers

import (
	db "Ecommerce/pkg/db"
	"Ecommerce/pkg/models"
	"Ecommerce/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var database *gorm.DB

func GetAllUsers() []models.User {
	database = db.GetDB()
	var users []models.User
	res := database.Find(&users)
	if res.Error != nil {
		log.Fatal(res.Error)
	}
	return users
}

func Signup(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid user data for model binding",
			"data":    err.Error(),
		})
		return
	}
	if user.Password == os.Getenv("ADMIN_PASS") && user.Email == os.Getenv("ADMIN_EMAIL") {
		user.UserType = "ADMIN"
	} else {
		user.UserType = "USER"
	}

	if user.UserType == "ADMIN" {
		var existingAdmin models.User
		result := database.Where("user_type = ?", "ADMIN").First(&existingAdmin)
		if result.Error == nil {
			c.JSON(400, gin.H{
				"status":  "error",
				"message": "Admin user already exists",
				"data":    result.Error,
			})
		}
	}

	// check if user already exists
	var existingUser models.User
	result := database.Where("email = ?", user.Email).First(&existingUser)
	if result.Error == nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "User already exists",
			"data":    result,
		})
		return
	}

	//hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to hash the password",
			"data":    err.Error(),
		})
		return
	}

	// Сохранение пользователя в базе данных
	user.Password = hashedPassword
	result = database.Create(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to insert a user",
			"data":    result.Error.Error(),
		})
		return
	}

	// sign jwt with user id and email
	signedToken, err := utils.CreateToken(strconv.Itoa(int(user.ID)), user.Email, user.UserType)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to create token",
			"data":    err.Error(),
		})
	}
	// add token to cookie session
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    signedToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &cookie)
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User signed up successfully",
		"data":    user,
	})
}
