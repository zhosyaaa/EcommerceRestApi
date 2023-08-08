package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func ValidateCredentialsMiddleware(c *gin.Context) {
	validate := validator.New()
	type UserInputCred struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	var user UserInputCred
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "missing field",
		})
		return
	}
	if err := validate.Struct(&user); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid request body",
			"err":     err.Error(),
		})
		return
	}
	c.Next()
}
