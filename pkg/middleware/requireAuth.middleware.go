package middleware

import (
	"Ecommerce/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RequireAuthMiddleware(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	token, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT token not found "})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Try to signin first",
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	id, email, userType, err := utils.VerifyToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token verification failed"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return

	}
	fmt.Println(id, email, userType, "in middleware")
	c.Set("id", id)
	c.Set("email", email)
	c.Set("userType", userType)
	c.Next()
}
