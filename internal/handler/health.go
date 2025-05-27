package handler

import (
	"gin-web-project/pkg/utils"
	"github.com/gin-gonic/gin"
)

func HealthCheck(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		utils.Success(c, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

}
