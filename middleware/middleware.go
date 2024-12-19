package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	jwt "washit-api/token"
)

func JWTAuth() gin.HandlerFunc {
	return JWT(jwt.AccessTokenType, "any")
}

func JTWAuthAdmin() gin.HandlerFunc {
	return JWT(jwt.AccessTokenType, "admin")
}

func JWTRefresh() gin.HandlerFunc {
	return JWT(jwt.RefreshTokenType, "any")
}

func JWT(tokenType string, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
			c.Abort()
			return
		}

		log.Println("Token: ", token)

		payload, err := jwt.ValidateToken(token)
		if err != nil || payload == nil || payload["type"] != tokenType {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		if role != payload["role"] && role != "any" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "role unauthorized"})
			c.Abort()
			return
		}

		c.Set("userId", payload["id"])
		c.Set("userRole", payload["role"])
		c.Next()
	}
}
