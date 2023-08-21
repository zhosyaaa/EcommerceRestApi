package controllers

import (
	"Ecommerce/pkg/config"
	"Ecommerce/pkg/models"
	interfaces "Ecommerce/pkg/repository/interface"
	"Ecommerce/pkg/utils"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type UserController struct {
	userService interfaces.UserRepository
}

func NewUserController(userService interfaces.UserRepository) *UserController {
	return &UserController{userService: userService}
}

func (s *UserController) Signup(c *gin.Context) {
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
	_, err := s.userService.GetByID(strconv.Itoa(int(user.ID)))
	if err == nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "User already exists",
		})
		return
	}

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

	err = s.userService.CreateUser(&user)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to insert a user",
			"data":    err,
		})
		return
	}

	signedToken, err := utils.CreateToken(strconv.Itoa(int(user.ID)), user.Email, user.UserType)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to create token ",
			"data":    err.Error(),
		})
		return
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    signedToken,
		Path:     "/api/v1",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &cookie)
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User signed up successfully",
		"data":    signedToken,
	})
}

func (s *UserController) Signin(c *gin.Context) {
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
			"message": "Invalid request body",
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
	password := req.Password
	email := req.Email

	existingUser, err := s.userService.GetByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
	if !utils.VerifyPassword(password, existingUser.Password) {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid credentials",
		})
		return
	}
	signedToken, err := utils.CreateToken(strconv.Itoa(int(existingUser.ID)), existingUser.Email, existingUser.UserType)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to create token",
			"data":    err.Error(),
		})
	}
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    signedToken,
		Path:     "/api/v1",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &cookie)

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User signed in successfully",
		"data":    existingUser,
	})
}

func (s *UserController) Signout(c *gin.Context) {
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/api/v1",
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

// /api/v1/users/auth/profile
func (s *UserController) Profile(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(401, gin.H{
			"status":  "error",
			"message": "User not authenticated",
		})
		return
	}

	user, err := s.userService.GetByID(id.(string))
	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "user does not exist",
			"data":    err,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User successfully found",
		"data":    user,
	})
}
