package controllers

import (
	"Ecommerce/pkg/db"
	"Ecommerce/pkg/models"
	"Ecommerce/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
	"time"
)

var database *gorm.DB

func Signup(c *gin.Context) {
	database = db.GetDB()
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

func Signin(c *gin.Context) {
	type SigninRequest struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	// parse request body
	var req SigninRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid request body",
			"data":    err.Error(),
		})
		return
	}
	// get email and password from request body
	password := req.Password
	email := req.Email

	// check if user already exists
	var existingUser models.User
	result := database.Where("email = ?", email).First(&existingUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"status":  "error",
				"message": "User does not exist",
			})
		} else {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "Internal server error",
			})
		}
		return
	}
	// check if password is correct
	if !utils.VerifyPassword(password, existingUser.Password) {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid credentials",
		})
		return
	}
	// sign jwt with user id and email
	signedToken, err := utils.CreateToken(strconv.Itoa(int(existingUser.ID)), existingUser.Email, existingUser.UserType)
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
		"data":    existingUser,
	})
}

func Signout(c *gin.Context) {
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &cookie)
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User logged out successfully",
		"data":    nil,
	})
}

func Profile(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	result := database.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"status":  "error",
				"message": "User does not exist",
			})
		} else {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "Internal server error",
			})
		}
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Successfully fetched user",
		"data":    user,
	})
}
