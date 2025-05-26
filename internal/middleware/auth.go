package middleware

import (
	"gin-web-project/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// 移除 "Bearer " 前缀
		if strings.HasPrefix(token, "Bearer ") {
			token = token[7:]
		}

		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
