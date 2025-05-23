package handler

import (
	"gin-web-project/pkg/utils"
	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	utils.Success(c, gin.H{
		"status":  "ok",
		"message": "Server is running",
	})
}
