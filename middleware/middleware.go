package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	jwt "washit-api/token"
)

func JWTAuth() gin.HandlerFunc {
	return JWT(jwt.AccessTokenType)
}

func JWTRefresh() gin.HandlerFunc {
	return JWT(jwt.RefreshTokenType)
}

func JWT(tokenType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		payload, err := jwt.ValidateToken(token)
		if err != nil || payload == nil || payload["type"] != tokenType {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		log.Println("Payload: ", payload["id"])
		c.Set("userId", payload["id"])
		c.Next()
	}
}
