package middleware

import (
	"Ecommerce/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func RequireAuthMiddleware(c *gin.Context) {
	var token string
	authHeader := c.GetHeader("Authorization")
	cookie, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT token not found "})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	fields := strings.Fields(authHeader)
	if len(fields) != 0 && fields[0] == "Bearer" {
		token = fields[1]
	} else if err == nil {
		token = cookie
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
	c.Set("id", id)
	c.Set("email", email)
	c.Set("userType", userType)
	c.Next()
}
