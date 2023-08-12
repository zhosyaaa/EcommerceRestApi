package controllers

import (
	"Ecommerce/pkg/config"
	"Ecommerce/pkg/db"
	"Ecommerce/pkg/models"
	"Ecommerce/pkg/utils"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

// /     /api/v1/users/auth/singup +
func Signup(c *gin.Context) {
	session := db.GetDB().Session(&gorm.Session{})
	validate := validator.New()
	request_id := c.GetString("x-request-id")
	logger := log.With().Str("request_id", request_id).Logger()
	logger.Debug().Msg("Received signup request")

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
		logger.Error().Err(err).Msg("Error occurred while binding request data")
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid user data for model binding",
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

	if user.Password == config.GetEnvVar("ADMIN_PASS") && user.Email == config.GetEnvVar("ADMIN_EMAIL") {
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
				"message": "Admin user already exists",
				"data":    result.Error,
			})
		}
	}

	var existingUser models.User
	result := session.Where("email = ?", user.Email).First(&existingUser)
	if result.Error == nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "User already exists",
			"data":    result,
		})
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		logger.Error().Err(err).Msg("Error occurred while hashing password")
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to hash the password",
			"data":    err.Error(),
		})
		return
	}
	user.Password = hashedPassword

	log.Info().Msg("Creating user")
	result = session.Create(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to insert a user",
			"data":    result.Error.Error(),
		})
		return
	}
	session.Commit()

	log.Info().Msg("Creating signed token")
	signedToken, err := utils.CreateToken(strconv.Itoa(int(user.ID)), user.Email, user.UserType)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to create token ",
			"data":    err.Error(),
		})
		return
	}

	log.Info().Msg("Setting JWT cookie")
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    signedToken,
		Path:     "/api/v1",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &cookie)

	log.Info().Str("request_id", request_id).Msg("Response sent")
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User signed up successfully",
		"data":    signedToken,
	})
}

// /api/v1/users/auth/singin +
func Signin(c *gin.Context) {
	session := db.GetDB().Session(&gorm.Session{})
	validate := validator.New()
	request_id := c.GetString("x-request-id")
	logger := log.With().Str("request_id", request_id).Logger()
	logger.Debug().Msg("Received signin request")

	type SigninRequest struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	var req SigninRequest
	if err := c.BindJSON(&req); err != nil {
		logger.Error().Err(err).Msg("Error occurred while binding request data")
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid request body",
			"data":    err.Error(),
		})
		return
	}
	if err := validate.Struct(&req); err != nil {
		logger.Error().Err(err).Msg("Invalid request body")
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid request body",
			"err":     err.Error(),
		})
		return
	}
	password := req.Password
	email := req.Email

	var existingUser models.User
	result := session.Where("email = ?", email).Preload("Address").Preload("UserCart").Preload("Orders").First(&existingUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.Error().Err(result.Error).Msg("User does not exist")
			c.JSON(404, gin.H{
				"status":  "error",
				"message": "User does not exist",
			})
		} else {
			logger.Error().Err(result.Error).Msg("Internal server error")
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "Internal server error",
			})
		}
		return
	}
	if !utils.VerifyPassword(password, existingUser.Password) {
		logger.Warn().Msg("Invalid credentials")
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid credentials",
		})
		return
	}

	logger.Info().Msg("Creating signed token")
	signedToken, err := utils.CreateToken(strconv.Itoa(int(existingUser.ID)), existingUser.Email, existingUser.UserType)
	if err != nil {
		log.Error().Err(err).Str("request_id", request_id).Msg("Failed to create token")
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to create token",
			"data":    err.Error(),
		})
	}

	log.Info().Msg("Setting JWT cookie")
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    signedToken,
		Path:     "/api/v1",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &cookie)

	log.Info().Msg("Response sent")
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User signed in successfully",
		"data":    existingUser,
	})
}

// /api/v1/users/auth/singout
func Signout(c *gin.Context) {
	request_id := c.GetString("x-request-id")
	logger := log.With().Str("request_id", request_id).Logger()
	logger.Debug().Msg("Received signout request")

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/api/v1",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &cookie)

	log.Info().Msg("User logged out successfully")
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User logged out successfully",
		"data":    nil,
	})
}

// /api/v1/users/auth/profile
func Profile(c *gin.Context) {
	session := db.GetDB().Session(&gorm.Session{})
	request_id := c.GetString("x-request-id")
	logger := log.With().Str("request_id", request_id).Logger()
	logger.Debug().Msg("Received profile request")

	var user models.User
	id, exists := c.Get("id")
	if !exists {
		logger.Warn().Msg("User not authenticated")
		c.JSON(401, gin.H{
			"status":  "error",
			"message": "User not authenticated",
		})
		return
	}
	logger.Debug().Msg("User authenticated")

	result := session.Where("ID=?", id).Preload("Address").Preload("UserCart").Preload("Orders").First(&user)
	if result.Error != nil {
		logger.Error().Err(result.Error).Msg("Failed to find user")
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "user does not exist",
			"data":    result.Error,
		})
		return
	}
	logger.Info().Int("user_id", int(user.ID)).Msg("User successfully found")
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User successfully found",
		"data":    user,
	})
}
