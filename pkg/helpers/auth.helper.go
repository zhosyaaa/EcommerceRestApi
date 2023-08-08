package helpers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckUserType(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("userType")
		if !exists || userType != role {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized to access this route",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func MatchUserTypeToUid(userId string) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, exists := c.Get("uid")
		userType, userTypeExists := c.Get("userType")
		if !exists || !userTypeExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized to access this route",
			})
			c.Abort()
			return
		}

		if userType == "USER" && uid != userId {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized to access this route",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
