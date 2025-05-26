package middleware

import (
	"gin-web-project/internal/database"
	"gin-web-project/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
			utils.Error(c, 401, "invalid token")
			c.Abort()
			return
		}
		redisToken := database.RedisClient.Get(database.Ctx, utils.TokenPrefix+strconv.Itoa(int(claims.UserID)))
		if redisToken.Val() == "" {
			utils.Error(c, 401, "Invalid token")
			c.Abort()
			return
		}
		//验证token是否在Redis中存在
		if redisToken.Val() != token {
			utils.Error(c, 401, "Invalid token")
		}

		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
