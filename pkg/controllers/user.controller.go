package controllers

import (
	"Ecommerce/pkg/db"
	"Ecommerce/pkg/models"
	"Ecommerce/pkg/utils"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

// /     /api/v1/users/auth/singup +
func Signup(c *gin.Context) {
	session := db.GetDB().Session(&gorm.Session{})
	validate := validator.New()

	type UserInputCred struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	var user models.User
	var userAddress models.Address
	userAddress.ZipCode = gofakeit.Zip()
	userAddress.City = gofakeit.City()
	userAddress.State = gofakeit.State()
	userAddress.Country = gofakeit.Country()
	userAddress.Street = gofakeit.Street()
	userAddress.HouseNumber = gofakeit.StreetNumber()
	user.Address = userAddress
	user.Orders = make([]models.Order, 0)
	user.UserCart = make([]models.ProductsToOrder, 0)

	var userBind UserInputCred
	if err := c.ShouldBindJSON(&userBind); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid user data for model binding Signup",
			"data":    err.Error(),
		})
		return
	}
	if err := validate.Struct(&userBind); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid request body",
			"err":     err.Error(),
		})
		return
	}
	user.Email = userBind.Email
	user.Username = userBind.Username
	user.Password = userBind.Password

	///для чтения с енв файла
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if user.Password == os.Getenv("ADMIN_PASS") && user.Email == os.Getenv("ADMIN_EMAIL") {
		user.UserType = "ADMIN"
	} else {
		user.UserType = "USER"
	}
	if user.UserType == "ADMIN" {
		var existingAdmin models.User
		result := session.Where("user_type = ?", "ADMIN").First(&existingAdmin)
		if result.Error == nil {
			c.JSON(400, gin.H{
				"status":  "error",
				"message": "Admin user already exists Signup",
				"data":    result.Error,
			})
		}
	}

	var existingUser models.User
	result := session.Where("email = ?", user.Email).First(&existingUser)
	if result.Error == nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "User already exists Signup",
			"data":    result,
		})
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to hash the password Signup",
			"data":    err.Error(),
		})
		return
	}
	user.Password = hashedPassword

	result = session.Create(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to insert a user Signup",
			"data":    result.Error.Error(),
		})
		return
	}
	session.Commit()

	signedToken, err := utils.CreateToken(user.ID, user.Email, user.UserType)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to create token Signup",
			"data":    err.Error(),
		})
		return
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    signedToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &cookie)
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User signed up successfully Signup",
		"data":    signedToken,
	})
}

// /api/v1/users/auth/singin +
func Signin(c *gin.Context) {
	session := db.GetDB().Session(&gorm.Session{})
	validate := validator.New()

	type SigninRequest struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	var req SigninRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid request body Signin",
			"data":    err.Error(),
		})
		return
	}
	if err := validate.Struct(&req); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid request body",
			"err":     err.Error(),
		})
		return
	}
	//username := req.Username
	password := req.Password
	email := req.Email

	var existingUser models.User
	result := session.Where("email = ?", email).First(&existingUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"status":  "error",
				"message": "User does not exist Signin",
			})
		} else {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "Internal server error Signin",
			})
		}
		return
	}
	if !utils.VerifyPassword(password, existingUser.Password) {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid credentials Signin",
		})
		return
	}

	signedToken, err := utils.CreateToken(existingUser.ID, existingUser.Email, existingUser.UserType)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to create token Signin",
			"data":    err.Error(),
		})
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    signedToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &cookie)
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User signed in successfully Signin",
		"data":    existingUser,
	})
}

// /api/v1/users/auth/singout
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
		"message": "User logged out successfully Signout",
		"data":    nil,
	})
}

// /api/v1/users/auth/profile
func Profile(c *gin.Context) {
	session := db.GetDB().Session(&gorm.Session{})
	var user models.User
	id, exists := c.Get("email")
	if !exists {
		c.JSON(401, gin.H{
			"status":  "error",
			"message": "User not authenticated Profile",
		})
		return
	}
	result := session.Where("email=?", id).First(&user)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "user does not exist",
			"data":    result.Error,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Successfully fetched user Profile",
		"data":    user,
	})
}
